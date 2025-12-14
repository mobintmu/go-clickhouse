package clickhouse

import (
	"context"
	"fmt"
	"go-clickhouse/internal/config"
	"go-clickhouse/internal/storage/clickhouse/product"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouse struct {
	conn    clickhouse.Conn
	Product *product.Repository
}

func (c *ClickHouse) Conn() clickhouse.Conn {
	return c.conn
}

func New(cfg *config.Config) (*ClickHouse, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.ClickHouse.Host + ":" + cfg.ClickHouse.Port},
		Auth: clickhouse.Auth{
			Database: cfg.ClickHouse.DB,
			Username: cfg.ClickHouse.User,
			Password: cfg.ClickHouse.Password,
		},
	})
	if err != nil {
		fmt.Println("🛑 Click House could not connect: ", err)
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		conn.Close()
		fmt.Println("🛑 Click House could not connect: ", err)
		return nil, err
	}

	return &ClickHouse{
		conn:    conn,
		Product: product.NewProductRepository(conn),
	}, nil
}
