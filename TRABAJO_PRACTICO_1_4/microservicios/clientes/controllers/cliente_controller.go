package controllers

import (
	"net/http"

	"clientes/models"
	"clientes/services"

	"github.com/gin-gonic/gin"
)

type ClienteController struct {
	Service *services.ClienteService
}

func (c *ClienteController) Crear(ctx *gin.Context) {
	var cliente models.Cliente
	if err := ctx.ShouldBindJSON(&cliente); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	creado, err := c.Service.CrearCliente(cliente)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, creado)
}

func (c *ClienteController) ObtenerPorID(ctx *gin.Context) {
	id := ctx.Param("id")

	cliente, err := c.Service.ObtenerCliente(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cliente)
}
