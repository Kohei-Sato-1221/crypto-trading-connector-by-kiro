package handler

import (
	"net/http"

	"github.com/crypto-trading-connector/backend/internal/generated"
	"github.com/crypto-trading-connector/backend/internal/service"
	"github.com/labstack/echo/v4"
)

// CryptoHandler handles HTTP requests for cryptocurrency endpoints
type CryptoHandler struct {
	service service.CryptoService
}

// NewCryptoHandler creates a new crypto handler
func NewCryptoHandler(service service.CryptoService) *CryptoHandler {
	return &CryptoHandler{
		service: service,
	}
}

// GetMarketData handles GET /api/v1/crypto/market
func (h *CryptoHandler) GetMarketData(c echo.Context) error {
	marketData, err := h.service.GetMarketData()
	if err != nil {
		return handleError(c, http.StatusInternalServerError, generated.INTERNALSERVERERROR, err.Error())
	}

	return c.JSON(http.StatusOK, marketData)
}

// GetCryptoByID handles GET /api/v1/crypto/:id
func (h *CryptoHandler) GetCryptoByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return handleError(c, http.StatusBadRequest, generated.BADREQUEST, "cryptocurrency ID is required")
	}

	cryptoData, err := h.service.GetCryptoByID(id)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "cryptocurrency not found: "+id {
			return handleError(c, http.StatusNotFound, generated.NOTFOUND, "Cryptocurrency not found")
		}
		return handleError(c, http.StatusInternalServerError, generated.INTERNALSERVERERROR, err.Error())
	}

	return c.JSON(http.StatusOK, cryptoData)
}

// GetChartData handles GET /api/v1/crypto/:id/chart
func (h *CryptoHandler) GetChartData(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return handleError(c, http.StatusBadRequest, generated.BADREQUEST, "cryptocurrency ID is required")
	}

	period := c.QueryParam("period")
	// Default period is handled in service layer

	chartData, err := h.service.GetChartData(id, period)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "cryptocurrency not found: "+id {
			return handleError(c, http.StatusNotFound, generated.NOTFOUND, "Cryptocurrency not found")
		}
		return handleError(c, http.StatusInternalServerError, generated.INTERNALSERVERERROR, err.Error())
	}

	return c.JSON(http.StatusOK, chartData)
}

// handleError is a helper function to return error responses
func handleError(c echo.Context, statusCode int, errorType generated.ErrorResponseError, message string) error {
	return c.JSON(statusCode, generated.ErrorResponse{
		Error:   errorType,
		Message: message,
	})
}
