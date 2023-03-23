package worker

import (
	"context"
	"log"
	"route256/libs/workerpool"
	"route256/loms/internal/domain"
	"route256/loms/internal/model"
	"time"
)

type DeleteReservationWorker struct {
	pool                *workerpool.Pool[model.JobDeleteReservation, bool]
	ordersRepository    domain.OrdersRepository
	warehouseRepository domain.WarehouseRepository
	ctx                 context.Context
}

func NewDeleteReservationWorker(
	ctx context.Context,
	amountWorkers int,
	ordersRepository domain.OrdersRepository,
	warehouseRepository domain.WarehouseRepository,
) *DeleteReservationWorker {
	DeleteReservationService := &DeleteReservationWorker{
		pool:                workerpool.NewPool[model.JobDeleteReservation, bool](ctx, amountWorkers),
		ordersRepository:    ordersRepository,
		warehouseRepository: warehouseRepository,
		ctx:                 ctx,
	}

	// запускаем прослушивание канала с результатом и логируем
	go func() {
		for result := range DeleteReservationService.pool.ResultChan {
			log.Printf("received result: %+v", result)
		}
	}()

	return DeleteReservationService
}

// AddDelayJob Добавление отложенной задачи на проверку оплаты заказа
func (w *DeleteReservationWorker) AddDelayJob(job model.JobDeleteReservation) {
	time.AfterFunc(time.Minute*10, func() {
		log.Printf("add job to JobChan: %+v", job.OrderId)
		w.pool.JobChan <- workerpool.Job[model.JobDeleteReservation, bool]{
			Callback: w.doWork,
			Args:     job,
		}
	})
}

// Выполняет проверку в каком статусе заказ, если он всё еще ожидает оплату, то отменяет его и
// удаляет резервы
func (w *DeleteReservationWorker) doWork(job model.JobDeleteReservation) (bool, error) {
	log.Printf("check order status %v", job.OrderId)
	order, err := w.ordersRepository.GetOrder(w.ctx, job.OrderId)
	if err != nil {
		return false, err
	}

	if order.Status == model.AwaitingPayment {
		err = w.ordersRepository.UpdateStatusOrder(w.ctx, job.OrderId, model.Cancelled)
		if err != nil {
			return false, err
		}

		log.Printf("order %v mark as %v", job.OrderId, model.Cancelled)

		err = w.warehouseRepository.DeleteReservation(w.ctx, job.OrderId)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}
