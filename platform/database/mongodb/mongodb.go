package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type MongoDB struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUsername string
	DBPassword string
}

var (
	client *mongo.Client
	db     *mongo.Database
)

// InitMongo initializes the MongoDB connection with advanced options
func InitMongo(config MongoDB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var uri string
	if config.DBUsername != "" && config.DBPassword != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s/%s", config.DBHost, config.DBPort, config.DBName)
	}

	clientOptions := options.Client().ApplyURI(uri).
		SetMaxPoolSize(100).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(24 * time.Hour).
		SetRetryWrites(true).
		SetRetryReads(true).
		SetWriteConcern(writeconcern.Majority()).
		SetReadPreference(readpref.SecondaryPreferred())

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Verify connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	db = client.Database(config.DBName)
	log.Println("[MONGO_DB] connected successfully")
	return nil
}

// GetCollection retrieves a collection from the database
func GetCollection(collectionName string) *mongo.Collection {
	if db == nil {
		log.Fatalf("Database not initialized")
	}
	return db.Collection(collectionName)
}

func isReplicaSet(client *mongo.Client) (bool, error) {
	var result bson.M
	err := client.Database("admin").RunCommand(context.Background(), bson.D{{Key: "isMaster", Value: 1}}).Decode(&result)
	if err != nil {
		return false, fmt.Errorf("failed to run isMaster command: %v", err)
	}

	if _, ok := result["setName"]; ok {
		return true, nil
	}

	return false, nil
}

// withTransaction handles operations with or without transactions based on the replica set status
func withTransaction(txnFn func(sessCtx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	isReplicaSet, err := isReplicaSet(client)
	if err != nil {
		return nil, err
	}

	if isReplicaSet {
		// If it's a replica set, use a session and transaction
		session, err := client.StartSession()
		if err != nil {
			return nil, fmt.Errorf("failed to start session: %v", err)
		}
		defer session.EndSession(context.Background())

		// Use the session and transaction
		result, err := session.WithTransaction(context.Background(), txnFn)
		if err != nil {
			return nil, fmt.Errorf("transaction failed: %v", err)
		}

		log.Println("[MONGO_DB] transaction successfully committed")
		return result, nil
	}

	// If it's not a replica set, execute the operation without a transaction
	log.Println("[MONGO_DB] not a replica set, running operation without transaction")

	// Create a fake session context by wrapping a regular context
	ctx := mongo.NewSessionContext(context.Background(), nil)
	return txnFn(ctx)
}

// FindAndUpdateWithRollback finds a document, updates it, and rolls back if there's an error
func FindAndUpdateWithRollback(collectionName string, filter interface{}, update interface{}) (interface{}, error) {
	return withTransaction(func(sessCtx mongo.SessionContext) (interface{}, error) {
		coll := GetCollection(collectionName)

		var updatedDoc bson.M
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true) // Return the updated document

		err := coll.FindOneAndUpdate(sessCtx, filter, update, opts).Decode(&updatedDoc)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, errors.New("no document was modified or found")
			}
			return nil, fmt.Errorf("failed to update document: %v", err)
		}

		return updatedDoc, nil
	})
}

// InsertDocumentWithRollback inserts a new document with rollback capability
func InsertDocumentWithRollback(collectionName string, document interface{}) (interface{}, error) {
	return withTransaction(func(sessCtx mongo.SessionContext) (interface{}, error) {
		coll := GetCollection(collectionName)

		result, err := coll.InsertOne(sessCtx, document)

		if err != nil {
			return nil, fmt.Errorf("failed to insert document: %v", err)
		}

		return result.InsertedID, nil
	})
}

// DeleteDocumentWithRollback deletes a document with rollback capability
func DeleteDocumentWithRollback(collectionName string, filter interface{}) (int64, error) {
	result, err := withTransaction(func(sessCtx mongo.SessionContext) (interface{}, error) {
		coll := GetCollection(collectionName)

		result, err := coll.DeleteOne(sessCtx, filter)
		if err != nil {
			return nil, fmt.Errorf("failed to delete document: %v", err)
		}

		return result.DeletedCount, nil
	})

	if err != nil {
		return 0, err
	}

	return result.(int64), nil
}

// FindDocuments finds documents in the specified collection based on a filter
func FindDocuments(collectionName string, filter interface{}, limit int64, sort interface{}) ([]bson.M, error) {
	coll := GetCollection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Handle nil filter
	var filterDoc interface{}
	if filter == nil {
		filterDoc = bson.M{}
	} else {
		filterDoc = filter
	}

	// Create find options
	findOptions := options.Find().SetLimit(limit)

	// Handle nil sort
	if sort != nil {
		findOptions.SetSort(sort)
	}

	cursor, err := coll.Find(ctx, filterDoc, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %v", err)
	}

	return results, nil
}

// FindOne finds a single document in the specified collection based on a filter
func FindOne(collectionName string, filter interface{}) (bson.M, error) {
	coll := GetCollection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result bson.M
	err := coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no document found matching the filter: %v", filter)
		}
		return nil, fmt.Errorf("failed to find document: %v", err)
	}

	return result, nil
}

// BulkWriteWithRollback performs multiple write operations with rollback capability
func BulkWriteWithRollback(collectionName string, operations []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	result, err := withTransaction(func(sessCtx mongo.SessionContext) (interface{}, error) {
		coll := GetCollection(collectionName)

		opts := options.BulkWrite().SetOrdered(false)
		result, err := coll.BulkWrite(sessCtx, operations, opts)
		if err != nil {
			return nil, fmt.Errorf("bulk write failed: %v", err)
		}

		return result, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*mongo.BulkWriteResult), nil
}

// AggregateWithRollback performs an aggregation pipeline with rollback capability
func AggregateWithRollback(collectionName string, pipeline interface{}) ([]bson.M, error) {
	result, err := withTransaction(func(sessCtx mongo.SessionContext) (interface{}, error) {
		coll := GetCollection(collectionName)

		cursor, err := coll.Aggregate(sessCtx, pipeline)
		if err != nil {
			return nil, fmt.Errorf("aggregation failed: %v", err)
		}
		defer cursor.Close(sessCtx)

		var results []bson.M
		if err = cursor.All(sessCtx, &results); err != nil {
			return nil, fmt.Errorf("failed to decode aggregation results: %v", err)
		}

		return results, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]bson.M), nil
}

// CreateIndexWithRollback creates an index with rollback capability
func CreateIndexWithRollback(collectionName string, keys bson.D, options *options.IndexOptions) (string, error) {
	result, err := withTransaction(func(sessCtx mongo.SessionContext) (interface{}, error) {
		coll := GetCollection(collectionName)

		indexName, err := coll.Indexes().CreateOne(sessCtx, mongo.IndexModel{Keys: keys, Options: options})
		if err != nil {
			return nil, fmt.Errorf("failed to create index: %v", err)
		}

		return indexName, nil
	})

	if err != nil {
		return "", err
	}

	return result.(string), nil
}

// DropIndexWithRollback drops an index with rollback capability
func DropIndexWithRollback(collectionName string, indexName string) error {
	_, err := withTransaction(func(sessCtx mongo.SessionContext) (interface{}, error) {
		coll := GetCollection(collectionName)

		_, err := coll.Indexes().DropOne(sessCtx, indexName)
		if err != nil {
			return nil, fmt.Errorf("failed to drop index: %v", err)
		}

		return nil, nil
	})

	return err
}

// Shutdown disconnects from the MongoDB database
func Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if client != nil {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatalf("Cannot disconnect MongoDB: %v", err)
		}
		log.Println("[MONGO_DB] disconnected successfully")
	}
}
