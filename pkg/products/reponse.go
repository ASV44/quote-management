package products

import (
	"quote-management-tech-task/db/sqlc"
	"quote-management-tech-task/models"
)

type GetProductsResponse struct {
	Products   []sqlc.Product    `json:"products"`
	Pagination models.Pagination `json:"pagination"`
}

type GetProductResponse struct {
	sqlc.Product
}
