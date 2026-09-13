package products

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	models "main/models/products"
	service "main/services/products"
)

type Handler struct{ service *service.Service }

func NewHandler(s *service.Service) *Handler { return &Handler{service: s} }

func (h *Handler) Register(r *gin.Engine) {
	r.GET("/products/search", h.search)
}

func (h *Handler) search(c *gin.Context) {
	text := c.Query("q")
	if text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "q es obligatorio"})
		return
	}

	limit := 10
	if raw := c.Query("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit debe ser un entero positivo"})
			return
		}
		limit = parsed
	}

	query := models.Query{
		Text:     text,
		Category: c.Query("category"),
		Brand:    c.Query("brand"),
		Limit:    limit,
	}

	results, err := h.service.Search(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}
