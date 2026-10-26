package repository

import (
	"context"
	"errors"

	"ecommerce-microservices/services/product/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	Create(context.Context, *model.Product) error
	GetByID(context.Context, int64) (*model.Product, error)
	GetAll(context.Context) ([]model.Product, error)
}

type PostgresProductRepository struct {
	db *pgxpool.Pool
}

func NewPostgresProductRepository(db *pgxpool.Pool) ProductRepository {
	return &PostgresProductRepository{db: db}
}

func (r *PostgresProductRepository) Create(ctx context.Context, p *model.Product) error {
	return r.db.QueryRow(ctx, `INSERT INTO products(name, price) VALUES($1, $2) RETURNING id`, p.Name, p.Price).Scan(&p.ID)
}

func (r *PostgresProductRepository) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	var p model.Product
	err := r.db.QueryRow(ctx, `SELECT id, name, price FROM products WHERE id=$1`, id).Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostgresProductRepository) GetAll(ctx context.Context) ([]model.Product, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, price FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []model.Product{}
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
