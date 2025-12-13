package storage

import (
	"go-clickhouse/internal/storage/cache"
	"go-clickhouse/internal/storage/clickhouse"
	"go-clickhouse/internal/storage/sql/sqlc"
)

type Storage struct {
	SQL        *sqlc.Queries
	Cache      *cache.Store
	ClickHouse *clickhouse.ClickHouse
}

func New(sql *sqlc.Queries, cache *cache.Store, ch *clickhouse.ClickHouse) *Storage {
	return &Storage{
		SQL:        sql,
		Cache:      cache,
		ClickHouse: ch,
	}
}
