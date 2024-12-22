package health

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

const okStatus = "OK"
const serviceName = "products-service"

// HealthHandler godoc
//
//	@Summary	Returns health status
//	@Tags		Other
//	@ID			health-status
//	@Produce	json
//	@Success	200	{object}	healthResponseBody
//	@Router		/api/health [get]
func HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, &healthResponseBody{
		Status:      okStatus,
		ServiceName: serviceName,
		Build:       os.Getenv("APP_BUILD"),
		DeployedAt:  os.Getenv("APP_DEPLOYED_AT"),
	})
}

type healthResponseBody struct {
	Status      string `json:"status" example:"OK"`
	ServiceName string `json:"serviceName" example:"identity"`
	Build       string `json:"build"`
	DeployedAt  string `json:"deployedAt"`
}
