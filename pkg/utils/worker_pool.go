package utils

import "sync"

type WorkerPool interface {
	Start()
	Add(work func())
	Wait()
	CloseTaskC()
}

// workerPool is a thread-safe implementation of the Worker-Pool concurrency pattern.
// Implemetes the WorkerPool interface
type workerPool struct {
	workerN   int
	taskQueue chan func()
	waitGroup *sync.WaitGroup
}

// NewWorkerPool initiates a new work pool
func NewWorkerPool(workerN, queueSize int) WorkerPool {
	queue := make(chan func(), queueSize)
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(workerN)
	return &workerPool{
		workerN:   workerN,
		taskQueue: queue,
		waitGroup: waitGroup,
	}
}

// Start runs `workerN` number of goroutines
// to process work added to `taskQueue`
func (p *workerPool) Start() {
	for i := 0; i < p.workerN; i++ {
		go func() {
			defer p.waitGroup.Done() // mark this goroutine as done
			for work := range p.taskQueue {
				work()
			}
		}()
	}
}

// Add adds work to the task queue
func (p *workerPool) Add(work func()) {
	p.taskQueue <- work
}

// CloseTaskC closes the task queue channel.
// It is the responsibility of the caller of `Add`
// to also call CloseTaskC.
func (p *workerPool) CloseTaskC() {
	close(p.taskQueue)
}

// Wait waits until all the worker goroutines have exited.
// It is the responsibility of the caller of `Exec` to also
// call `Wait`. This would otherwise lead to *leaking goroutines*.
func (p *workerPool) Wait() {
	p.waitGroup.Wait()
}
