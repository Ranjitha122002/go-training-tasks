package repository

import (
	"context"
	"errors"
	"time"

	"ecommerce-microservices/services/order/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	Create(context.Context, *model.Order) error
	GetByID(context.Context, int64) (*model.Order, error)
	ListAll(context.Context) ([]*model.Order, error)
	UpdateStatus(context.Context, int64, model.OrderStatus) (*model.Order, error)
}

type PostgresOrderRepository struct {
	db *pgxpool.Pool
}

func NewPostgresOrderRepository(db *pgxpool.Pool) OrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) Create(ctx context.Context, o *model.Order) error {
	if o.CreatedAt.IsZero() {
		o.CreatedAt = time.Now()
	}

	return r.db.QueryRow(ctx,
		`INSERT INTO orders(customer_id, product_id, quantity, total, status, created_at) 
		 VALUES($1, $2, $3, $4, $5, $6) RETURNING id`,
		o.CustomerID, o.ProductID, o.Quantity, o.Total, o.Status, o.CreatedAt).Scan(&o.ID)
}

func (r *PostgresOrderRepository) GetByID(ctx context.Context, id int64) (*model.Order, error) {
	var o model.Order
	err := r.db.QueryRow(ctx,
		`SELECT id, customer_id, product_id, quantity, total, status, created_at 
		 FROM orders WHERE id=$1`, id).Scan(
		&o.ID, &o.CustomerID, &o.ProductID, &o.Quantity, &o.Total, &o.Status, &o.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *PostgresOrderRepository) ListAll(ctx context.Context) ([]*model.Order, error) {
	rows, err := r.db.Query(ctx, 
		`SELECT id, customer_id, product_id, quantity, total, status, created_at 
		 FROM orders ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.ProductID, &o.Quantity, &o.Total, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}
	return orders, rows.Err()
}

func (r *PostgresOrderRepository) UpdateStatus(ctx context.Context, id int64, status model.OrderStatus) (*model.Order, error) {
	var o model.Order
	err := r.db.QueryRow(ctx,
		`UPDATE orders SET status = $1 WHERE id = $2 
		 RETURNING id, customer_id, product_id, quantity, total, status, created_at`,
		status, id).Scan(&o.ID, &o.CustomerID, &o.ProductID, &o.Quantity, &o.Total, &o.Status, &o.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &o, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
