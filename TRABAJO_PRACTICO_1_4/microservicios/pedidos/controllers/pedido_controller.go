package controllers

import (
	"net/http"

	"pedidos/services"

	"github.com/gin-gonic/gin"
)

type PedidoController struct {
	Service *services.PedidoService
}

func (c *PedidoController) ListarProductos(ctx *gin.Context) {
	productos, err := c.Service.ListarProductos()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"productos": productos})
}

type confirmarPedidoRequest struct {
	ClienteID  string `json:"cliente_id"`
	ProductoID string `json:"producto_id"`
}

func (c *PedidoController) ConfirmarPedido(ctx *gin.Context) {
	var req confirmarPedidoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pedido, err := c.Service.ConfirmarPedido(req.ClienteID, req.ProductoID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"mensaje":   "pedido confirmado",
		"pedido_id": pedido.ID,
	})
}
