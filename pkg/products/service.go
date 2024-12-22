package products

import (
	"context"
	"fmt"
	"math"

	"quote-management-tech-task/db/sqlc"
	"quote-management-tech-task/models"
)

type Service struct {
	db *sqlc.Queries
}

func NewService(db *sqlc.Queries) Service {
	return Service{
		db: db,
	}
}

func (s Service) CreateProducts(ctx context.Context, products []CreateProductData) error {
	createProductsParams := make([]sqlc.CreateProductsParams, 0, len(products))
	for _, product := range products {
		createProductsParams = append(createProductsParams, sqlc.CreateProductsParams{
			Name:        product.Name,
			Description: &product.Description,
			Price:       product.Price,
			TaxRate:     product.TaxRate,
			Metadata:    product.Metadata,
		})
	}

	if err := s.db.CreateProducts(ctx, createProductsParams).Close(); err != nil {
		return fmt.Errorf("failed to bulk insert products: %w", err)
	}

	return nil
}

func (s Service) GetProducts(
	ctx context.Context,
	queryParams GetProductsQueries,
) (GetProductsResponse, error) {
	var sortBy string
	// Convert payload param to column name and add default in case it is empty or not supported or invalid value
	switch queryParams.SortBy {
	case Name, Price:
		sortBy = string(queryParams.SortBy)
	case CreatedAt:
		sortBy = "created_at"
	default:
		sortBy = string(Name)
	}

	// Convert payload param to sorting order and add default in case it is empty or invalid
	var sortOrder string
	switch queryParams.SortOrder {
	case Ascending, Descending:
		sortOrder = string(queryParams.SortOrder)
	default:
		sortOrder = string(Ascending)
	}

	const defaultPerPage = 10
	if queryParams.PerPage == 0 || queryParams.PerPage < 0 {
		queryParams.PerPage = defaultPerPage
	}

	if queryParams.Page <= 0 {
		queryParams.Page = 1
	}

	products, err := s.db.GetProducts(ctx, sqlc.GetProductsParams{
		Query:     queryParams.FilterQuery,
		Sortby:    sortBy,
		Sortorder: sortOrder,
		Limit:     queryParams.PerPage,
		Offset:    (queryParams.Page - 1) * queryParams.PerPage,
	})
	if err != nil {
		return GetProductsResponse{}, err
	}

	totalCount, err := s.db.GetProductsTotalCount(ctx, queryParams.FilterQuery)
	if err != nil {
		return GetProductsResponse{}, err
	}

	return GetProductsResponse{
		Products: products,
		Pagination: models.Pagination{
			PerPage:    int(queryParams.PerPage),
			Page:       int(queryParams.Page),
			TotalPages: int(math.Ceil(float64(totalCount) / float64(queryParams.PerPage))),
			TotalCount: int(totalCount),
		},
	}, nil
}

func (s Service) GetProductByID(ctx context.Context, productID int32) (GetProductResponse, error) {
	product, err := s.db.GetProductByID(ctx, productID)
	if err != nil {
		return GetProductResponse{}, err
	}

	return GetProductResponse{Product: product}, nil
}

func (s Service) UpdateProducts(ctx context.Context, products []UpdateProductData) error {
	updateProductsParams := make([]sqlc.UpdateProductsParams, 0, len(products))
	for _, product := range products {
		updateProductsParams = append(updateProductsParams, sqlc.UpdateProductsParams{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			TaxRate:     product.TaxRate,
			Metadata:    product.Metadata,
		})
	}

	if err := s.db.UpdateProducts(ctx, updateProductsParams).Close(); err != nil {
		return fmt.Errorf("failed to bulk insert products: %w", err)
	}

	return nil
}

func (s Service) DeleteProducts(ctx context.Context, productIDs []int32) error {
	if err := s.db.DeleteProducts(ctx, productIDs); err != nil {
		return fmt.Errorf("failed to bulk delete products: %w", err)
	}

	return nil
}
