package main

import (
	"log"

	"clientes/controllers"
	"clientes/repositories"
	"clientes/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Capa de repositorio: datos en memoria
	repo := repositories.NewClientesMemRepo()

	// Capa de servicio: lógica de negocio
	service := &services.ClienteService{Repo: repo}

	// Capa de controlador: HTTP
	controller := &controllers.ClienteController{Service: service}

	router := gin.Default()
	router.POST("/clientes", controller.Crear)
	router.GET("/clientes/:id", controller.ObtenerPorID)

	log.Println("Microservicio clientes escuchando en http://localhost:8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
