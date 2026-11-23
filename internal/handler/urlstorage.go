package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntduc6403/bookmark-management/internal/service"
)

type ShortenURLRequest struct {
	ExpTime int    `json:"exp"`
	URL     string `json:"url" binding:"required,url"`
}

type ShortenURLResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ShortenURLHandler struct {
	service service.ShortenUrl
}

func NewShortenURLHandler(svc service.ShortenUrl) *ShortenURLHandler {
	return &ShortenURLHandler{
		service: svc,
	}
}

// Shorten handles shortening a long URL into a short code.
// @Summary Shorten URL
// @Description Shortens a given URL with an optional expiration time in seconds.
// @Tags Links
// @Accept json
// @Produce json
// @Param request body ShortenURLRequest true "Shorten URL request payload"
// @Success 200 {object} ShortenURLResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/links/shorten [post]
func (h *ShortenURLHandler) Shorten(c *gin.Context) {
	var req ShortenURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	if req.ExpTime < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Expiration time must not be negative",
		})
		return
	}
	code, err := h.service.ShortenUrl(
		c,
		req.URL,
		req.ExpTime,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to shorten URL",
		})
		return
	}

	c.JSON(http.StatusOK, ShortenURLResponse{
		Code:    code,
		Message: "Shorten URL generated successfully!",
	})
}
