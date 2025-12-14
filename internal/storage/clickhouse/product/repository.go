package product

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type Product struct {
	ID          int32
	Name        string
	Description string
	Price       int64
}

type Repository struct {
	conn clickhouse.Conn
}

func NewProductRepository(conn clickhouse.Conn) *Repository {
	return &Repository{conn: conn}
}

const productColumns = `id, name, description, price`

func (r *Repository) SelectProduct(ctx context.Context, id int32) (*Product, error) {
	var p Product

	err := r.conn.QueryRow(ctx, `
        SELECT `+productColumns+`
        FROM products
        WHERE id = ?
        LIMIT 1
    `, id).Scan(&p.ID, &p.Name, &p.Description, &p.Price)

	if err != nil {
		return nil, err
	}

	return &p, nil
}
