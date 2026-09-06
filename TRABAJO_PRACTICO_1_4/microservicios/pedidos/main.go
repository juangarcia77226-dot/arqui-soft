package main

import (
	"log"
	"time"

	"pedidos/controllers"
	"pedidos/messaging"
	"pedidos/repositories"
	"pedidos/services"

	"github.com/gin-gonic/gin"
)

const amqpURI = "amqp://user:pass@localhost:5672"
const cacheTTL = 30 * time.Second

func main() {
	// Publisher de RabbitMQ: declara la cola y queda listo para publicar
	// el evento pedido.confirmado.
	publisher, err := messaging.NewRabbitMQPublisher(amqpURI)
	if err != nil {
		log.Fatalf("no se pudo conectar a RabbitMQ: %v", err)
	}
	defer publisher.Close()

	// Capa de repositorio: productos (con caché) y pedidos, en memoria.
	productosRepo := repositories.NewProductosMemRepo()
	productosCacheados := repositories.NewProductosCachedRepo(productosRepo, cacheTTL)
	pedidosRepo := repositories.NewPedidosMemRepo()

	// Capa de servicio: lógica de negocio + publicación del evento.
	service := &services.PedidoService{
		ProductosRepo: productosCacheados,
		PedidosRepo:   pedidosRepo,
		Publisher:     publisher,
	}

	// Capa de controlador: HTTP.
	controller := &controllers.PedidoController{Service: service}

	router := gin.Default()
	router.GET("/productos", controller.ListarProductos)
	router.POST("/pedidos", controller.ConfirmarPedido)

	log.Println("Microservicio pedidos escuchando en http://localhost:8082")
	if err := router.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}
