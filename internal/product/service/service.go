package service

import (
	"context"
	"go-clickhouse/internal/config"
	"go-clickhouse/internal/product/dto"
	"go-clickhouse/internal/storage"
	"go-clickhouse/internal/storage/sql/sqlc"

	"go.uber.org/zap"
)

type Product struct {
	store *storage.Storage
	log   *zap.Logger
	cfg   *config.Config
}

func New(s *storage.Storage,
	log *zap.Logger,
	cfg *config.Config) *Product {
	return &Product{
		store: s,
		log:   log,
		cfg:   cfg,
	}
}

func (s *Product) Create(ctx context.Context, req dto.AdminCreateProductRequest) (dto.ProductResponse, error) {
	arg := sqlc.CreateProductParams{
		ProductName:        req.Name,
		ProductDescription: req.Description,
		Price:              req.Price,
		IsActive:           true,
	}
	product, err := s.store.SQL.CreateProduct(ctx, arg)
	if err != nil {
		return dto.ProductResponse{}, err
	}
	s.log.Info("Product created", zap.Int32("id", product.ID))
	s.store.Cache.Set(ctx, s.store.Cache.KeyProduct(product.ID), product, s.cfg.Redis.DefaultTTL)
	s.store.Cache.Delete(ctx, s.store.Cache.KeyAllProducts())
	return dto.ProductResponse{
		ID:          product.ID,
		Name:        product.ProductName,
		Description: product.ProductDescription,
		Price:       product.Price,
	}, nil
}

func (s *Product) Update(ctx context.Context, req dto.AdminUpdateProductRequest) (dto.ProductResponse, error) {
	arg := sqlc.UpdateProductParams{
		ID:                 int32(req.ID),
		ProductName:        req.Name,
		ProductDescription: req.Description,
		Price:              req.Price,
		IsActive:           req.IsActive,
	}
	product, err := s.store.SQL.UpdateProduct(ctx, arg)
	if err != nil {
		return dto.ProductResponse{}, err
	}
	s.store.Cache.Set(ctx, s.store.Cache.KeyProduct(product.ID), product, s.cfg.Redis.DefaultTTL)
	s.store.Cache.Delete(ctx, s.store.Cache.KeyAllProducts())
	return dto.ProductResponse{
		ID:          product.ID,
		Name:        product.ProductName,
		Description: product.ProductDescription,
		Price:       product.Price,
	}, nil
}

func (s *Product) Delete(ctx context.Context, id int32) error {
	s.store.Cache.Delete(ctx, s.store.Cache.KeyProduct(id))
	s.store.Cache.Delete(ctx, s.store.Cache.KeyAllProducts())
	return s.store.SQL.DeleteProduct(ctx, id)
}

func (s *Product) GetProductByID(ctx context.Context, id int32) (dto.ProductResponse, error) {
	var product sqlc.Product
	err := s.store.Cache.Get(ctx, s.store.Cache.KeyProduct(id), &product)
	if err != nil {
		product, err = s.store.SQL.GetProduct(ctx, id)
		if err != nil {
			return dto.ProductResponse{}, err
		}
		s.store.Cache.Set(ctx, s.store.Cache.KeyProduct(product.ID), product, s.cfg.Redis.DefaultTTL)
	}
	result := dto.ProductResponse{
		ID:          product.ID,
		Name:        product.ProductName,
		Description: product.ProductDescription,
		Price:       product.Price,
	}
	return result, nil
}

func (s *Product) ListProducts(ctx context.Context) (dto.ClientListProductsResponse, error) {
	var resp []dto.ProductResponse
	if err := s.store.Cache.Get(ctx, s.store.Cache.KeyAllProducts(), &resp); err == nil {
		return resp, nil
	}
	products, err := s.store.SQL.ListProducts(ctx)
	if err != nil {
		return nil, err
	}
	resp = make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		resp = append(resp, dto.ProductResponse{
			ID:          product.ID,
			Name:        product.ProductName,
			Description: product.ProductDescription,
			Price:       product.Price,
		})
	}
	s.store.Cache.Set(ctx, s.store.Cache.KeyAllProducts(), resp, s.cfg.Redis.DefaultTTL)
	return resp, nil
}

func (s *Product) GetProductWithReport(ctx context.Context, id int32) (dto.ProductResponse, error) {
	productCH, err := s.store.ClickHouse.Product.SelectProduct(ctx, id)
	if err != nil {
		return dto.ProductResponse{}, err
	}
	return dto.ProductResponse{
		ID:          productCH.ID,
		Name:        productCH.Name,
		Description: productCH.Description,
		Price:       productCH.Price,
	}, nil
}
