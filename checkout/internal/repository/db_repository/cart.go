package db_repository

import (
	"context"
	"errors"
	"log"
	"route256/checkout/internal/model"
	"route256/checkout/internal/repository/db_repository/schema"
	"route256/checkout/internal/repository/db_repository/transactor"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/multierr"
)

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i CartRepository -o ./mocks/ -s "_minimock.go"

type CartRepository interface {
	AddToCart(ctx context.Context, userID int64, sku uint32, count uint32) error
	UpdateCountCart(ctx context.Context, userID int64, sku uint32, count uint32) error
	DeleteCart(ctx context.Context, userID int64, sku uint32) error
	DeleteUserCart(ctx context.Context, userID int64) error
	ListCart(ctx context.Context, userID int64) ([]model.CartItem, error)
	GetCartItem(ctx context.Context, userID int64, sku uint32) (model.CartItem, error)
}

type cartRepository struct {
	transactor.QueryEngineProvider
}

func NewCartRepository(provider transactor.QueryEngineProvider) *cartRepository {
	return &cartRepository{
		provider,
	}
}

var (
	ErrNotFound       = errors.New("not found")
	ErrCartRepository = errors.New("cart repository error")
)

const cartTable = "cart"

func (r *cartRepository) AddToCart(ctx context.Context, userID int64, sku uint32, count uint32) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Insert(cartTable).
		Columns(
			"user_id",
			"sku",
			"count",
		).
		Values(userID, sku, count).
		Suffix("ON CONFLICT ON CONSTRAINT cart_pk DO UPDATE SET count = EXCLUDED.count").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return multierr.Append(ErrCartRepository, err)
	}

	log.Printf("AddToCart Query: %v with %v", query, args)
	_, err = db.Exec(ctx, query, args...)
	if err != nil {
		return multierr.Append(ErrCartRepository, err)
	}
	return nil
}

func (r *cartRepository) UpdateCountCart(ctx context.Context, userID int64, sku uint32, count uint32) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Update(cartTable).
		Set("count", count).
		Where(sq.Eq{"user_id": userID}).
		Where(sq.Eq{"sku": sku}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return multierr.Append(ErrCartRepository, err)
	}

	log.Printf("UpdateCountCart Query: %v with %v", query, args)
	_, err = db.Exec(ctx, query, args...)
	if err != nil {
		return multierr.Append(ErrCartRepository, err)
	}

	return nil
}

func (r *cartRepository) DeleteCart(ctx context.Context, userID int64, sku uint32) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Delete(cartTable).
		Where(sq.Eq{"user_id": userID}).
		Where(sq.Eq{"sku": sku}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return multierr.Append(ErrCartRepository, err)
	}

	_, err = db.Exec(ctx, query, args...)
	if err != nil {
		return multierr.Append(ErrCartRepository, err)
	}

	return nil
}

func (r *cartRepository) DeleteUserCart(ctx context.Context, userID int64) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Delete(cartTable).
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return multierr.Append(ErrCartRepository, err)
	}

	_, err = db.Exec(ctx, query, args...)
	if err != nil {
		return multierr.Append(ErrCartRepository, err)
	}

	return nil
}

func (r *cartRepository) ListCart(ctx context.Context, userID int64) ([]model.CartItem, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Select(
			"sku",
			"count",
		).
		From(cartTable).
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, multierr.Append(ErrCartRepository, err)
	}

	var item schema.Cart
	var result []model.CartItem

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, multierr.Append(ErrCartRepository, err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&item.Sku, &item.Count)
		if err != nil {
			return nil, err
		}
		result = append(result, model.CartItem{
			Sku:   uint32(item.Sku.Int),
			Count: uint32(item.Count.Int),
		})
	}

	return result, nil
}

func (r *cartRepository) GetCartItem(ctx context.Context, userID int64, sku uint32) (model.CartItem, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Select(
			"sku",
			"count",
		).
		From(cartTable).
		Where(sq.Eq{"user_id": userID}).
		Where(sq.Eq{"sku": sku}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return model.CartItem{}, multierr.Append(ErrCartRepository, err)
	}

	var item schema.Cart

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return model.CartItem{}, multierr.Append(ErrCartRepository, err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&item.Sku, &item.Count)
		if err != nil {
			return model.CartItem{}, err
		}
		return model.CartItem{
			Sku:   uint32(item.Sku.Int),
			Count: uint32(item.Count.Int),
		}, nil
	}

	return model.CartItem{}, multierr.Append(ErrCartRepository, ErrNotFound)
}
