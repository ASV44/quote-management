package products

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
