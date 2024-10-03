package workerpool

import (
	"container/heap"
	"context"
	"fmt"
	"sync"
	"time"
)

// Task represents a unit of work with priority and timeout.
type Task struct {
	Job        func(ctx context.Context) error
	Ctx        context.Context
	CancelFunc context.CancelFunc
	Priority   int
	Timeout    time.Duration
	CreatedAt  time.Time
}

// NewTask creates a new Task with context, priority, and timeout.
func NewTask(job func(ctx context.Context) error, priority int, timeout time.Duration) *Task {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &Task{
		Job:        job,
		Ctx:        ctx,
		CancelFunc: cancel,
		Priority:   priority,
		Timeout:    timeout,
		CreatedAt:  time.Now(),
	}
}

// priorityQueue implements heap.Interface and holds Tasks.
type priorityQueue []*Task

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x interface{}) {
	item := x.(*Task)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// WorkerPool manages the worker pool, task queue, and dynamic scaling.
type WorkerPool struct {
	tasks           priorityQueue
	taskLock        sync.Mutex
	workerWg        sync.WaitGroup
	maxWorkers      int
	minWorkers      int
	activeWorkers   int
	mu              sync.RWMutex
	scalingInterval time.Duration
	ctx             context.Context
	cancelFunc      context.CancelFunc
	metrics         *Metrics
	shutdownCh      chan struct{}
	workerSem       chan struct{}
	lastScaleTime   time.Time
	cooldownPeriod  time.Duration
	tasksAdded      int
	tasksCompleted  int

	taskMap map[string]struct{} // Map to track task uniqueness
}

// Metrics holds statistics about the worker pool's performance.
type Metrics struct {
	TasksCompleted        int64
	TasksFailed           int64
	AverageWaitTime       time.Duration
	AverageProcessingTime time.Duration
	mu                    sync.Mutex
}

// MetricsSnapshot represents a snapshot of the metrics without the mutex.
type MetricsSnapshot struct {
	TasksCompleted        int64
	TasksFailed           int64
	AverageWaitTime       time.Duration
	AverageProcessingTime time.Duration
}

// NewWorkerPool creates a new WorkerPool with advanced features.
func NewWorkerPool(ctx context.Context, minWorkers, maxWorkers int, scalingInterval time.Duration) *WorkerPool {
	ctx, cancel := context.WithCancel(ctx)
	return &WorkerPool{
		tasks:           make(priorityQueue, 0),
		minWorkers:      minWorkers,
		maxWorkers:      maxWorkers,
		activeWorkers:   minWorkers,
		scalingInterval: scalingInterval,
		ctx:             ctx,
		cancelFunc:      cancel,
		metrics:         &Metrics{},
		shutdownCh:      make(chan struct{}),
		workerSem:       make(chan struct{}, maxWorkers),
		lastScaleTime:   time.Now(),
		cooldownPeriod:  5 * time.Second,
		taskMap:         make(map[string]struct{}), // Initialize task map
	}
}

// AddTask adds a task to the worker pool's priority queue.
func (wp *WorkerPool) AddTask(task *Task) {
	taskID := fmt.Sprintf("%p", task.Job) // Example task ID based on the function pointer. Customize as needed.

	wp.taskLock.Lock()
	defer wp.taskLock.Unlock()

	// Check for duplicates
	if _, exists := wp.taskMap[taskID]; exists {
		fmt.Printf("Duplicate task detected. Task not added.\n")
		return
	}

	// Add task to the map to track it
	wp.taskMap[taskID] = struct{}{}

	heap.Push(&wp.tasks, task)
	wp.tasksAdded++
	fmt.Printf("Task added. Total tasks: %d\n", wp.tasksAdded)
}

// CancelTask cancels a specific task.
func (wp *WorkerPool) CancelTask(task *Task) {
	task.CancelFunc()
}

// worker is a goroutine that processes tasks.
func (wp *WorkerPool) worker(id int) {
	defer wp.workerWg.Done()
	for {
		select {
		case <-wp.ctx.Done():
			return
		case <-wp.shutdownCh:
			return
		default:
			wp.processNextTask(id)
		}
	}
}

// processNextTask handles the next task in the queue.
func (wp *WorkerPool) processNextTask(id int) {
	wp.taskLock.Lock()
	if wp.tasks.Len() == 0 {
		wp.taskLock.Unlock()
		time.Sleep(100 * time.Millisecond)
		return
	}
	task := heap.Pop(&wp.tasks).(*Task)
	taskID := fmt.Sprintf("%p", task.Job) // Use the same task ID logic
	wp.taskLock.Unlock()

	startTime := time.Now()
	wp.metrics.updateWaitTime(startTime.Sub(task.CreatedAt))

	select {
	case <-task.Ctx.Done():
		if task.Ctx.Err() == context.DeadlineExceeded {
			fmt.Printf("Worker %d: Task timed out\n", id)
		} else {
			fmt.Printf("Worker %d: Task canceled\n", id)
		}
		wp.metrics.incrementTasksFailed()
	default:
		if err := task.Job(task.Ctx); err != nil {
			fmt.Printf("Worker %d: Error: %v\n", id, err)
			wp.metrics.incrementTasksFailed()
		} else {
			wp.metrics.incrementTasksCompleted()
		}
	}

	wp.metrics.updateProcessingTime(time.Since(startTime))

	// Critical section to update and print task completion safely
	wp.mu.Lock()
	wp.tasksCompleted++
	completed := wp.tasksCompleted
	added := wp.tasksAdded

	// Clean up task from taskMap
	wp.taskLock.Lock()
	delete(wp.taskMap, taskID)
	wp.taskLock.Unlock()

	wp.mu.Unlock()

	fmt.Printf("Task completed. Completed tasks: %d/%d\n", completed, added)
}

// Run starts the initial workers and the auto-scaling process.
func (wp *WorkerPool) Run() {
	for i := 0; i < wp.minWorkers; i++ {
		wp.workerWg.Add(1)
		wp.workerSem <- struct{}{} // Add this line
		go wp.worker(i)
	}

	go wp.autoScale()
}

// autoScale automatically scales the number of workers based on the task queue load.
func (wp *WorkerPool) autoScale() {
	ticker := time.NewTicker(wp.scalingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case <-ticker.C:
			wp.mu.Lock()
			taskQueueSize := wp.tasks.Len()
			now := time.Now()
			if now.Sub(wp.lastScaleTime) > wp.cooldownPeriod {
				if taskQueueSize > wp.activeWorkers && wp.activeWorkers < wp.maxWorkers {
					wp.scaleUp(taskQueueSize)
					wp.lastScaleTime = now
				} else if taskQueueSize < wp.activeWorkers/2 && wp.activeWorkers > wp.minWorkers {
					wp.scaleDown()
					wp.lastScaleTime = now
				}
			}
			wp.mu.Unlock()
		}
	}
}

// scaleUp increases the number of workers.
func (wp *WorkerPool) scaleUp(taskQueueSize int) {
	newWorkers := min(taskQueueSize-wp.activeWorkers, wp.maxWorkers-wp.activeWorkers)
	for i := 0; i < newWorkers; i++ {
		wp.workerSem <- struct{}{}
		wp.workerWg.Add(1) // Add this line
		go wp.worker(wp.activeWorkers + i)
	}
	wp.activeWorkers += newWorkers
	fmt.Printf("Scaled up to %d workers\n", wp.activeWorkers)
}

// scaleDown reduces the number of workers.
func (wp *WorkerPool) scaleDown() {
	workersToRemove := (wp.activeWorkers - wp.minWorkers) / 2 // Remove half of the excess workers
	if workersToRemove == 0 {
		workersToRemove = 1 // Remove at least one worker
	}
	for i := 0; i < workersToRemove; i++ {
		<-wp.workerSem
	}
	wp.activeWorkers -= workersToRemove
	fmt.Printf("Scaled down to %d workers\n", wp.activeWorkers)
}

// Wait waits for all tasks to complete.
// Wait waits for all tasks to complete.
func (wp *WorkerPool) Wait() {
	fmt.Println("Waiting for all tasks to complete...")

	// Wait for tasks to be completed
	for {
		wp.mu.RLock()
		completed := wp.tasksCompleted
		added := wp.tasksAdded
		wp.mu.RUnlock()

		if completed >= added {
			break
		}

		fmt.Printf("Still waiting... Completed tasks: %d/%d\n", completed, added)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("All tasks completed.")
}

// Shutdown gracefully shuts down the worker pool.
func (wp *WorkerPool) Shutdown() {
	fmt.Println("Shutting down worker pool...")
	wp.cancelFunc()
	close(wp.shutdownCh)
	wp.workerWg.Wait()
	wp.taskLock.Lock()
	wp.tasks = priorityQueue{}
	wp.taskLock.Unlock()
	fmt.Println("Worker pool shut down.")
}

// GetMetrics returns a snapshot of the current metrics of the worker pool.
func (wp *WorkerPool) GetMetrics() MetricsSnapshot {
	wp.metrics.mu.Lock()
	defer wp.metrics.mu.Unlock()
	return MetricsSnapshot{
		TasksCompleted:        wp.metrics.TasksCompleted,
		TasksFailed:           wp.metrics.TasksFailed,
		AverageWaitTime:       wp.metrics.AverageWaitTime,
		AverageProcessingTime: wp.metrics.AverageProcessingTime,
	}
}

// Helper functions for Metrics
func (m *Metrics) incrementTasksCompleted() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TasksCompleted++
}

func (m *Metrics) incrementTasksFailed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TasksFailed++
}

func (m *Metrics) updateWaitTime(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.AverageWaitTime = (m.AverageWaitTime*time.Duration(m.TasksCompleted) + duration) / time.Duration(m.TasksCompleted+1)
}

func (m *Metrics) updateProcessingTime(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.AverageProcessingTime = (m.AverageProcessingTime*time.Duration(m.TasksCompleted) + duration) / time.Duration(m.TasksCompleted+1)
}

// Helper function to calculate the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Test() {
	// Create a worker pool with 4 to 10 workers, with scaling check every 5 seconds.
	ctx := context.Background()
	pool := NewWorkerPool(ctx, 4, 10, 5*time.Second)

	// Start the worker pool.
	pool.Run()

	// Add some tasks to the pool.
	for i := 0; i < 20; i++ {
		taskNum := i
		priority := 1 //rand.Intn(5) // Random priority between 0 and 4
		task := NewTask(func(ctx context.Context) error {
			fmt.Printf("Processing task %d with priority %d\n", taskNum, priority)
			time.Sleep(1 * time.Second)
			return nil
		}, priority, 5*time.Second) // 5 second timeout for each task
		pool.AddTask(task)
	}

	// Wait for all tasks to complete.
	pool.Wait()

	// Add a small delay to ensure all goroutines have completed
	time.Sleep(1 * time.Second)

	// Get and print metrics
	metrics := pool.GetMetrics()
	fmt.Printf("Tasks completed: %d\n", metrics.TasksCompleted)
	fmt.Printf("Tasks failed: %d\n", metrics.TasksFailed)
	fmt.Printf("Average wait time: %v\n", metrics.AverageWaitTime)
	fmt.Printf("Average processing time: %v\n", metrics.AverageProcessingTime)

	// Shut down the pool after all tasks are done.
	pool.Shutdown()
}
