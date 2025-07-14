package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type (
	// WorkerParty manages a group of workers,
	// allowing dynamic scaling via deltaCh,
	// sending jobs via jobCh, and clean shutdown of all workers.
	WorkerParty struct {
		jobCh   chan Job // Channel through which workers receive jobs to execute.
		deltaCh chan int // Channel to signal desired worker count for dynamic scaling.
		done    bool     // Flag indicating if the WorkerParty is shutting down.

		workers []*worker // Slice of currently active workers.

		mu  sync.Mutex      // Mutex protecting concurrent access to WorkerParty state.
		wg  sync.WaitGroup  // WaitGroup to track worker completion before shutdown.
		ctx context.Context // Context controlling lifecycle (cancellation, timeout) for all workers.
	}

	// worker represents a single worker that listens for jobs on jobCh
	// and executes them until receiving a stop signal or context cancellation.
	worker struct {
		id   int           // Unique identifier for the worker.
		stop chan struct{} // Channel signaling the worker to stop.
	}

	// Job defines a unit of work that a worker will execute.
	Job func()
)

func NewWorker(id int) *worker {
	return &worker{
		id:   id,
		stop: make(chan struct{}, 1),
	}
}

func (w *worker) Start(ctx context.Context, wg *sync.WaitGroup, jobCh chan Job) {
	defer wg.Done()

	for {
		select {
		case <-w.stop:
			fmt.Println("Завершаем работу воркера с помощью stop: ", w.id)
			return
		case task := <-jobCh:
			task()
			fmt.Println("Воркер сделал задачу: ", w.id)

		case <-ctx.Done():
			fmt.Println("Завершаем работу воркера с помощью контекста: ", w.id)
			return
		}
	}
}

func NewWorkerParty() *WorkerParty {
	return &WorkerParty{
		workers: make([]*worker, 0),
		deltaCh: make(chan int, 1),
		jobCh:   make(chan Job, 1),
		done:    true,
	}
}

// Передавать контекст который будет прослушиваться
func (wp *WorkerParty) Start(ctx context.Context) {
	wp.mu.Lock()
	wp.ctx = ctx
	wp.mu.Unlock()

	for {
		select {
		// Обновляем кол-во воркеров
		case delta := <-wp.deltaCh:
			wp.mu.Lock()
			current := len(wp.workers)
			wp.mu.Unlock()

			if delta > current {
				wp.AddWorkers(delta - current)
			} else {
				wp.RemoveWorkers(current - delta)
			}
			fmt.Printf("Number of workers changed from %d to %d\n", current, delta)
		case <-wp.ctx.Done():
			wp.vacation()
			return
		}
	}

}

func (wp *WorkerParty) Stop() {
	wp.done = true
	wp.wg.Wait()
}

func (wp *WorkerParty) Scale(workerCount int) error {
	/// валидировать count
	if workerCount < 0 {
		return errors.New("worker count must be >= 0")
	}

	wp.mu.Lock()
	done := wp.done
	wp.mu.Unlock()

	if !done {
		return errors.New("worker pool is inactive")
	}

	wp.deltaCh <- workerCount
	return nil
}

func (wp *WorkerParty) vacation() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	fmt.Println("Отправляем воркеров в отпуск...")
	for _, worker := range wp.workers {
		worker.stop <- struct{}{}
		close(worker.stop)
	}
	fmt.Println("Все воркеры в отпуске...")
}

func (wp *WorkerParty) AddWorkers(n int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	startId := len(wp.workers)
	for i := 1; i <= n; i++ {
		w := NewWorker(startId + i)

		wp.workers = append(wp.workers, w)

		wp.wg.Add(1)
		go w.Start(wp.ctx, &wp.wg, wp.jobCh)
		fmt.Println("Добавлен воркер:", w.id)
	}
	fmt.Println("Текущее количество воркеров:", len(wp.workers))
}

func (wp *WorkerParty) RemoveWorkers(n int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	for i := 0; i < n; i++ {
		idx := len(wp.workers) - 1

		w := wp.workers[idx]
		wp.workers = wp.workers[:idx]

		w.stop <- struct{}{}
		close(w.stop)
		fmt.Println("Удален воркер:", w.id)
	}
	fmt.Println("Текущее количество воркеров:", len(wp.workers))
}
