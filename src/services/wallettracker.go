package services

import (
	"context"
	"fmt"
	"pnl-scan-tool/package/workerpool"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/gofiber/fiber/v2"
	_ "github.com/swaggo/fiber-swagger" // Fiber Swagger middleware
)

type WalletTrackerRequest struct {
	WalletAddress string `json:"walletaddress" example:"EBw6beJFQePbH1x9WzMX5ipBBr634drKX2N1bCzJVDwY"`
}

// TaskResponse represents the response for a successful task creation
type WalletTrackerResponse struct {
	Message  string `json:"message"`
	TaskID   string `json:"taskId"`
	Priority int    `json:"priority"`
	Timeout  string `json:"timeout"`
}

// ErrorResponse represents the response for an error
type ErrorResponse struct {
	Error string `json:"error"`
}

type WalletTrackerTaskManager struct {
	pool    *workerpool.WorkerPool
	taskMap map[string]*workerpool.Task // Task ID -> Task
	mu      sync.Mutex
}

func NewWalletTrackerTaskManager(minWorkers, maxWorkers int, scalingInterval time.Duration) *WalletTrackerTaskManager {
	pool := workerpool.NewWorkerPool(context.Background(), minWorkers, maxWorkers, scalingInterval)
	pool.Run()
	return &WalletTrackerTaskManager{
		pool:    pool,
		taskMap: make(map[string]*workerpool.Task),
	}
}

// AddWalletTrackerHandler creates a new task and adds it to the worker pool
// @Summary Add a new task to the task manager
// @Description Creates a task with a specified duration, priority, and timeout, and adds it to the worker pool
// @Tags add wallet tracker
// @Accept json
// @Produce json
// @Param task body WalletTrackerRequest true "Task Details"
// @Success 201 {object} WalletTrackerResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/wallettracker/add [post]
func (tm *WalletTrackerTaskManager) AddWalletTrackerHandler(ctx context.Context, telegramBot *bot.Bot) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tm.mu.Lock()
		defer tm.mu.Unlock()

		var walletTracker WalletTrackerRequest

		if err := c.BodyParser(&walletTracker); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error: "Failed to parse request body",
			})
		}

		wallet := walletTracker.WalletAddress
		priority := 1               // Random priority between 0 and 4
		timeout := 10 * time.Second // Timeout after 10 seconds

		// Create a new task
		task := workerpool.NewTask(func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				fmt.Printf("Task %s was cancelled or timed out\n", wallet)
				return ctx.Err()
			default:

				// var transfer solscan.Transfer

				// fmt.Printf("Processing tracking wallet %s with priority %d\n", wallet, priority)

				// solscan := solscan.Solscan{
				// 	Address:      wallet,
				// 	ExcludeToken: "So11111111111111111111111111111111111111111",
				// 	Flow:         "in",
				// }

				for {
					// transactions, err := solscan.GetTransactionsWallet()

					// transfer = transactions[0]

					// fmt.Println(transfer)

					fmt.Println(time.Now().Unix())

					_, err := telegramBot.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: -1001973873960,
						Text:   "test",
					})

					//time.Sleep(1 * time.Second)

					//telegramBot.SendMessage(fmt.Sprintf("%s: %s", wallet, transfer))

					if err != nil {
						fmt.Println("Err:", err)
						fmt.Printf("Finished %s\n", wallet)
						return nil
					}
				}
			}
		}, priority, timeout)

		// Add task to pool and map
		tm.taskMap[wallet] = task
		tm.pool.AddTask(task)

		return c.Status(fiber.StatusCreated).JSON(WalletTrackerResponse{
			Message:  "Wallet tracker added",
			TaskID:   wallet,
			Priority: priority,
			Timeout:  timeout.String(),
		})
	}
}

func (tm *WalletTrackerTaskManager) CancelTaskHandler(c *fiber.Ctx) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	taskID := c.Query("taskId")
	task, exists := tm.taskMap[taskID]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	// Cancel the task
	tm.pool.CancelTask(task)
	delete(tm.taskMap, taskID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task canceled",
		"taskId":  taskID,
	})
}

func (tm *WalletTrackerTaskManager) ListTasksHandler(c *fiber.Ctx) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	taskList := make([]fiber.Map, 0, len(tm.taskMap))
	for taskID, task := range tm.taskMap {
		taskList = append(taskList, fiber.Map{
			"taskId":   taskID,
			"priority": task.Priority,
			"timeout":  task.Timeout.String(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tasks": taskList,
	})
}

func (tm *WalletTrackerTaskManager) GetMetricsHandler(c *fiber.Ctx) error {
	metrics := tm.pool.GetMetrics()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tasksCompleted":        metrics.TasksCompleted,
		"tasksFailed":           metrics.TasksFailed,
		"averageWaitTime":       metrics.AverageWaitTime.String(),
		"averageProcessingTime": metrics.AverageProcessingTime.String(),
	})
}

func (tm *WalletTrackerTaskManager) ShutdownHandler(c *fiber.Ctx) error {
	tm.pool.Shutdown()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Worker pool is shutting down",
	})
}
