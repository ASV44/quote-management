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

const defaultPerPage = 10

func (s Service) GetProducts(
	ctx context.Context,
	queryParams GetProductsQueries,
) (GetProductsResponse, error) {
	var sortBy string
	// Convert payload param to column name
	switch queryParams.SortBy {
	case Name:
		sortBy = "name"
	case Price:
		sortBy = "price"
	case CreatedAt:
		sortBy = "created_at"
	default:
		sortBy = "name"
	}

	if queryParams.SortOrder == Descending {
		sortBy = fmt.Sprintf("%s DESC", sortBy)
	}

	if queryParams.PerPage == 0 || queryParams.PerPage < 0 {
		queryParams.PerPage = defaultPerPage
	}

	if queryParams.Page <= 0 {
		queryParams.Page = 1
	}

	products, err := s.db.GetProducts(ctx, sqlc.GetProductsParams{
		Query:  queryParams.FilterQuery,
		Sortby: sortBy,
		Limit:  queryParams.PerPage,
		Offset: (queryParams.Page - 1) * queryParams.PerPage,
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
