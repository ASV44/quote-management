package products

import (
	"encoding/json"

	"quote-management-tech-task/types"
)

type CreateProductsRequest struct {
	Products []CreateProductData `json:"products"`
}

type CreateProductData struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	TaxRate     float64         `json:"taxRate"`
	Metadata    json.RawMessage `json:"metadata"`
}

type ProductsSortBy string

const (
	Name      ProductsSortBy = "name"
	Price     ProductsSortBy = "price"
	CreatedAt ProductsSortBy = "creation date"
)

type SortOrder string

const (
	Ascending  SortOrder = "asc"
	Descending SortOrder = "desc"
)

type GetProductsQueries struct {
	FilterQuery string         `query:"query"`
	SortBy      ProductsSortBy `query:"sortBy" example:"[optional]desc"`
	SortOrder   SortOrder      `query:"sortOrder" example:"[optional]desc"`
	Page        int32          `query:"page"`
	PerPage     int32          `query:"perPage"`
}

type UpdateProductsRequest struct {
	Products []UpdateProductData `json:"products"`
}

type UpdateProductData struct {
	ID          int32             `json:"id"`
	Name        *string           `json:"name"`
	Description *string           `json:"description"`
	Price       types.NullFloat64 `json:"price"`
	TaxRate     types.NullFloat64 `json:"taxRate"`
	Metadata    json.RawMessage   `json:"metadata"`
}

type DeleteProductsRequest struct {
	ProductIDs []int32 `json:"productIDs"`
}
