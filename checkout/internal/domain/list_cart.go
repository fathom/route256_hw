package domain

import (
	"context"
	"fmt"
	"log"
	"route256/checkout/internal/config"
	"route256/checkout/internal/logger"
	"route256/checkout/internal/model"
	"route256/libs/workerpool"
	"sync"

	"github.com/pkg/errors"
)

func (d *domain) ListCart(ctx context.Context, userID int64) ([]model.CartItem, error) {
	logger.Debug(fmt.Sprintf("listCart for user: %+v", userID))

	userCart, err := d.cartRepository.ListCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Колбэк для обработки задания
	callback := func(cartItem model.CartItem) (model.CartItem, error) {
		logger.Debug(fmt.Sprintf("gorutine GetProduct start: %+v", cartItem.Sku))

		err := d.limiter.Wait(ctx)
		if err != nil {
			return model.CartItem{}, err
		}

		logger.Debug(fmt.Sprintf("gorutine GetProduct make request: %+v", cartItem.Sku))
		cartItem.Name, cartItem.Price, err = d.productService.GetProduct(ctx, cartItem.Sku)
		if err != nil {
			logger.Debug(fmt.Sprintf("get error from productService: %+v", err))
			return model.CartItem{}, errors.WithMessage(err, "wrong sku")
		}
		logger.Debug(fmt.Sprintf("gorutine GetProduct finish: %+v", cartItem.Sku))
		return cartItem, nil
	}

	pool, err := workerpool.NewPool[model.CartItem, model.CartItem](ctx, config.ConfigData.CountWorkers)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	// Добавление заданий в канал pool.JobChan
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			logger.Debug("closed gorutine add to JobChan")
		}()

		logger.Debug(fmt.Sprintf("total jobs %+v", len(userCart)))
		for _, cartItem := range userCart {
			select {
			case <-ctx.Done():
				logger.Debug("while add job context Done and close JobChan")
				close(pool.JobChan)
				return
			default:
				logger.Debug(fmt.Sprintf("add job: %+v", cartItem.Sku))
				pool.AddJob(workerpool.Job[model.CartItem, model.CartItem]{
					Callback: callback,
					Args:     cartItem,
				})
			}
		}
		logger.Debug("close chain JobChan")
		close(pool.JobChan)
	}()

	outputCart := make([]model.CartItem, 0, len(userCart))

	// Ожидаем результатов от воркеров
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			logger.Debug("closed gorutine read ResultChan")
		}()

		logger.Debug("start wait result from ResultChan")
		for cartItem := range pool.ResultChan {
			logger.Debug(fmt.Sprintf("received result: %+v", cartItem.Sku))
			outputCart = append(outputCart, cartItem)
		}
	}()

	// Ждем ошибок если они будут
	for err := range pool.ErrorChan {
		if err != nil {
			return nil, err
		}
	}

	log.Printf("")
	logger.Debug("wait finish work in ListCart")
	wg.Wait()

	logger.Debug("return result from ListCart")
	return outputCart, nil
}
