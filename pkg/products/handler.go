package products

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"quote-management-tech-task/models"
)

type Products interface {
	GetProducts(ctx context.Context, params GetProductsQueries) (GetProductsResponse, error)
	GetProductByID(ctx context.Context, productID int32) (GetProductResponse, error)
}

// Handler handles HTTP requests for specific route
type Handler struct {
	service Products
}

func NewHandler(service Products) Handler {
	return Handler{
		service: service,
	}
}

func (h Handler) Register(route *echo.Group) {
	product := route.Group("/products")
	product.GET("/", h.GetProducts)
	product.POST("/:productID", h.GetProduct)
}

// GetProducts godoc
//
//	@Summary	Products listing. Filterable and sortable via query params. Paginated.
//	@Tags		Products Listing
//	@ID			products-listing
//	@Param		query	query		string						false	"Query string for filtering and lookup name and description. Optional."
//	@Param		page	query		number						false	"Number of results per page. Optional, defaults to 10."
//	@Param		perPage	query		number						false	"Page number. Optional, defaults to 0."
//	@Success	200		{object}	GetProductsResponse	"Products listing"
//	@Failure	400		{object}	models.ErrorResponse
//	@Failure	500
//	@Router		/v1/products [get]
func (h Handler) GetProducts(c echo.Context) error {
	var queryParams GetProductsQueries
	if err := c.Bind(&queryParams); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "failed to bind query parameters"})
	}

	productResponse, err := h.service.GetProducts(c.Request().Context(), queryParams)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, productResponse)
}

// GetProduct godoc
//
//	@Summary	Product Details.
//	@Tags		product
//	@ID			product-details
//	@Param		productID	path		string			true	"Id of the product"
//	@Success	200					{object}	GetProductResponse	"Product details"
//	@Failure	404					"Product not found"
//	@Failure	500
//	@Router		/v1/products/{productID} [get]
func (h Handler) GetProduct(c echo.Context) error {
	productIDParam := c.Param("productID")

	productID, err := strconv.Atoi(productIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "invalid product id"})
	}

	productDetail, err := h.service.GetProductByID(c.Request().Context(), int32(productID))
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return c.NoContent(http.StatusNotFound)
	case err != nil:
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, productDetail)
}
