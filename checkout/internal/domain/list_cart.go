package domain

import (
	"context"
	"log"
	"route256/checkout/internal/config"
	"route256/checkout/internal/model"
	"route256/libs/workerpool"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

func (d *domain) ListCart(ctx context.Context, userID int64) ([]model.CartItem, error) {
	log.Printf("listCart for user: %+v", userID)

	userCart, err := d.cartRepository.ListCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Ограничиваем кол-во запросов 10rps
	limiter := rate.NewLimiter(rate.Every(1*time.Second/10), 10)

	// Колбэк для обработки задания
	callback := func(cartItem model.CartItem) (model.CartItem, error) {
		log.Printf("gorutine GetProduct start: %+v", cartItem.Sku)
		err := limiter.Wait(ctx)
		if err != nil {
			return model.CartItem{}, err
		}

		log.Printf("gorutine GetProduct make request: %+v", cartItem.Sku)
		cartItem.Name, cartItem.Price, err = d.productService.GetProduct(ctx, cartItem.Sku)
		if err != nil {
			return model.CartItem{}, errors.WithMessage(err, "wrong sku")
		}
		log.Printf("gorutine GetProduct finish: %+v", cartItem.Sku)
		return cartItem, nil
	}

	pool := workerpool.NewPool[model.CartItem, model.CartItem](ctx, config.ConfigData.CountWorkers)

	var wg sync.WaitGroup

	// Добавление заданий в канал pool.JobChan
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			log.Printf("closed gorutine add to JobChan")
		}()
		log.Printf("total jobs %+v", len(userCart))
		for _, cartItem := range userCart {
			select {
			case <-ctx.Done():
				return
			default:
				log.Printf("add job: %+v", cartItem.Sku)
				pool.AddJob(workerpool.Job[model.CartItem, model.CartItem]{
					Callback: callback,
					Args:     cartItem,
				})
			}
		}
		log.Printf("close chain JobChan")
		close(pool.JobChan)
	}()

	outputCart := make([]model.CartItem, 0, len(userCart))

	// Ожидаем результатов от воркеров
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			log.Printf("closed gorutine read ResultChan")
		}()

		log.Printf("start wait result from ResultChan")
		for cartItem := range pool.ResultChan {
			log.Printf("received result: %+v", cartItem.Sku)
			outputCart = append(outputCart, cartItem)
		}
	}()

	// Ждем ошибок если они будут
	for err := range pool.ErrorChan {
		if err != nil {
			return nil, err
		}
	}

	log.Printf("wait finish work in ListCart")
	wg.Wait()

	log.Printf("return result from ListCart")
	return outputCart, nil
}
