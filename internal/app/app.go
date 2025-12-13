package app

import (
	"go-clickhouse/internal/config"
	"go-clickhouse/internal/health"
	"go-clickhouse/internal/pkg/logger"
	productController "go-clickhouse/internal/product/controller"
	productService "go-clickhouse/internal/product/service"
	"go-clickhouse/internal/server"
	"go-clickhouse/internal/storage/cache"
	"go-clickhouse/internal/storage/sql"
	"go-clickhouse/internal/storage/sql/migrate"
	"go-clickhouse/internal/storage/sql/sqlc"

	"go.uber.org/fx"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(
			logger.NewLogger,
			config.NewConfig,
			sql.InitialDB,
			//server
			health.New,
			server.NewGinEngine,
			server.CreateHTTPServer,
			server.CreateGRPCServer,
			//db
			migrate.NewRunner, // migration runner
			sqlc.New,
			//cache
			cache.NewClient,
			cache.NewCacheStore,
			//controller
			productController.NewAdmin,
			productController.NewClient,
			productController.NewGRPC,
			//service
			productService.New,
		),
		fx.Invoke(
			server.RegisterRoutes,
			server.StartHTTPServer,
			server.StartGRPCServer,
			//migration
			migrate.RunMigrations,
			//life cycle
			logger.RegisterLoggerLifecycle,
			server.GRPCLifeCycle,
		),
	)
}
