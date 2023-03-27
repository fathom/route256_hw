package workerpool

import (
	"context"
	"errors"
	"log"
	"sync"
)

// Job Задача для выполенения в воркере
type Job[In, Out any] struct {
	Callback func(In) (Out, error)
	Args     In
}

type WorkerPool[In, Out any] interface {
	AddJob(Job[In, Out])
	Close()
}

// Pool Имплементация интерфейса WorkerPool
type Pool[In, Out any] struct {
	amountWorkers int // Кол-во воркеров в пуле

	JobChan    chan Job[In, Out] // Канал для отправки в него задач для воркеров
	ResultChan chan Out          // Канал для отправки результатов работы воркеров
	ErrorChan  chan error        // Канал для получения ошибки в воркерах

	wg sync.WaitGroup
}

// Проверка на соблюдение интерфейса WorkerPool
var _ WorkerPool[any, any] = &Pool[any, any]{}
var ErrNoWorkers = errors.New("no workers")

// NewPool создать новый worker pool
func NewPool[In, Out any](ctx context.Context, amountWorkers int) (*Pool[In, Out], error) {
	if amountWorkers < 1 {
		return nil, ErrNoWorkers
	}

	pool := &Pool[In, Out]{
		amountWorkers: amountWorkers,
		JobChan:       make(chan Job[In, Out], amountWorkers),
		ResultChan:    make(chan Out, amountWorkers),
		ErrorChan:     make(chan error, amountWorkers),
	}

	log.Printf("total workers %v", amountWorkers)
	for i := 0; i < amountWorkers; i++ {
		pool.wg.Add(1)
		go func(num int) {
			defer func() {
				pool.wg.Done()
				log.Printf("destroy worker %v", num)
			}()

			worker(num, ctx, pool.JobChan, pool.ResultChan, pool.ErrorChan)
		}(i)
	}

	go func() {
		pool.Close()
	}()

	return pool, nil
}

// AddJob добавить задачу в воркер пул
func (p *Pool[In, Out]) AddJob(job Job[In, Out]) {
	p.JobChan <- job
}

// Close закрыть каналы
func (p *Pool[In, Out]) Close() {
	p.wg.Wait()
	log.Printf("close chain ErrorChan and ResultChan")
	close(p.ErrorChan)
	close(p.ResultChan)
}

// worker для обработки задач
func worker[In, Out any](
	num int,
	ctx context.Context,
	jobChan <-chan Job[In, Out],
	resultChan chan<- Out,
	errorChan chan<- error,
) {
	log.Printf("add worker %v %v", num, jobChan)
	for job := range jobChan {
		log.Printf("listen jobChan")
		select {
		case <-ctx.Done():
			log.Printf("while lisen jobChan context Done")
			return
		default:
			out, err := job.Callback(job.Args)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- out
		}
	}
	log.Printf("exit from worker %v", num)
}
