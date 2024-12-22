package products

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"quote-management-tech-task/models"
)

type Products interface {
	CreateProducts(ctx context.Context, products []CreateProductData) error
	GetProducts(ctx context.Context, params GetProductsQueries) (GetProductsResponse, error)
	GetProductByID(ctx context.Context, productID int32) (GetProductResponse, error)
	UpdateProducts(ctx context.Context, products []UpdateProductData) error
	DeleteProducts(ctx context.Context, productIDs []int32) error
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
	route.POST("/products", h.CreateProducts)
	route.GET("/products", h.GetProducts)
	route.GET("/products/:productID", h.GetProduct)
	route.PUT("/products", h.UpdateProducts)
	route.DELETE("/products", h.DeleteProducts)
}

// CreateProducts godoc
//
//	@Summary	Create products. Bulk operation.
//	@Tags		Create Products
//	@ID			create-products
//	@Param		body	body	CreateProductsRequest	false	"create products data"
//	@Success	201
//	@Failure	400	{object}	models.ErrorResponse
//	@Failure	500
//	@Router		/v1/products [post]
func (h Handler) CreateProducts(c echo.Context) error {
	var createProductsReq CreateProductsRequest
	if err := c.Bind(&createProductsReq); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "failed to get create products request payload"})
	}

	if err := h.service.CreateProducts(c.Request().Context(), createProductsReq.Products); err != nil {
		log.Errorf("failed to create products: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

// GetProducts godoc
//
//	@Summary	Products listing. Filterable and sortable via query params. Paginated.
//	@Tags		Products Listing
//	@ID			products-listing
//	@Param		query		query		string				false	"Query string for filtering and lookup name and description. Optional."
//	@Param		sortBy		query		string				false	"Sorting column option. Optional. Defaults to 'name'"
//	@Param		sortOrder	query		string				false	"Sorting order. Either asc or desc. Optional, defaults to asc."
//	@Param		page		query		number				false	"Number of results per page. Optional, defaults to 10."
//	@Param		perPage		query		number				false	"Page number. Optional, defaults to 0."
//	@Success	200			{object}	GetProductsResponse	"Products listing"
//	@Failure	400			{object}	models.ErrorResponse
//	@Failure	500
//	@Router		/v1/products [get]
func (h Handler) GetProducts(c echo.Context) error {
	var queryParams GetProductsQueries
	if err := c.Bind(&queryParams); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "failed to bind query parameters"})
	}

	productResponse, err := h.service.GetProducts(c.Request().Context(), queryParams)
	if err != nil {
		log.Errorf("failed to get products: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, productResponse)
}

// GetProduct godoc
//
//	@Summary	Product Details.
//	@Tags		Product Details
//	@ID			product-details
//	@Param		productID	path		string				true	"Id of the product"
//	@Success	200			{object}	GetProductResponse	"Product details"
//	@Failure	404			"Product not found"
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
		log.Errorf("failed to get product by ID: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, productDetail)
}

// UpdateProducts godoc
//
//	@Summary	Update products details. Bulk operation
//	@Tags		Update products
//	@ID			update-products
//	@Param		body	body	UpdateProductsRequest	false	"update product details data"
//	@Success	204
//	@Failure	400	{object}	models.ErrorResponse
//	@Failure	500
//	@Router		/v1/products/{productID} [put]
func (h Handler) UpdateProducts(c echo.Context) error {
	var updateProductsReq UpdateProductsRequest
	if err := c.Bind(&updateProductsReq); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "failed to get update products request payload"})
	}

	if err := h.service.UpdateProducts(c.Request().Context(), updateProductsReq.Products); err != nil {
		log.Errorf("failed to update products: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteProducts godoc
//
//	@Summary	Delete products. Bulk operation
//	@Tags		Delete products
//	@ID			delete-products
//	@Param		body	body	DeleteProductsRequest	false	"Delete product IDs"
//	@Success	204
//	@Failure	400	{object}	models.ErrorResponse
//	@Failure	500
//	@Router		/v1/products/{productID} [delete]
func (h Handler) DeleteProducts(c echo.Context) error {
	var deleteProductsReq DeleteProductsRequest
	if err := c.Bind(&deleteProductsReq); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "failed to get delete products request payload"})
	}

	if err := h.service.DeleteProducts(c.Request().Context(), deleteProductsReq.ProductIDs); err != nil {
		log.Errorf("failed to delete products: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}
