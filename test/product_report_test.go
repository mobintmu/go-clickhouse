package test

import (
	"encoding/json"
	"fmt"
	"go-clickhouse/internal/config"
	"go-clickhouse/internal/product/dto"
	"io"
	"net/http"
	"testing"
)

func TestProductReportClient(t *testing.T) {
	// t.Parallel()
	WithHttpTestServer(t, func() {
		cfg, err := config.NewConfig()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		addr := fmt.Sprintf("http://%s:%d", cfg.HTTPAddress, cfg.HTTPPort)

		// the product has been create in different way
		product := dto.ProductResponse{
			ID:          1,
			Name:        "Test Product for Client Report",
			Description: "A product created for testing client report endpoint",
			Price:       1999,
		}
		clientGetProductByIDWithReport(t, product, addr)
	})
}

func clientGetProductByIDWithReport(t *testing.T, product dto.ProductResponse, addr string) {
	t.Run("List Products (Client)", func(t *testing.T) {
		resp, err := http.Get(addr + "/api/v1/products/" + fmt.Sprintf("%d/report", product.ID))
		if err != nil {
			t.Fatalf("Failed to send GET request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
			responseBody, _ := io.ReadAll(resp.Body)
			t.Logf(ResponseBodyMessage, string(responseBody))
			return
		}

		var result dto.ProductResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf(FailedToDecodeMessage, err)
		}

		if result.ID != product.ID {
			t.Fatal("Product not found in client list")
		}
	})
}
