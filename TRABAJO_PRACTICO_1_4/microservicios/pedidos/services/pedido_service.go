package services

import (
	"errors"

	"pedidos/messaging"
	"pedidos/models"
	"pedidos/repositories"
)

type PedidoService struct {
	ProductosRepo repositories.ProductosRepo
	PedidosRepo   repositories.PedidosRepo
	Publisher     *messaging.RabbitMQPublisher
}

func (s *PedidoService) ListarProductos() ([]models.Producto, error) {
	return s.ProductosRepo.ListarProductos()
}

func (s *PedidoService) ConfirmarPedido(clienteID, productoID string) (models.Pedido, error) {
	if clienteID == "" || productoID == "" {
		return models.Pedido{}, errors.New("cliente_id y producto_id son obligatorios")
	}

	pedido, err := s.PedidosRepo.CrearPedido(models.Pedido{
		ClienteID:  clienteID,
		ProductoID: productoID,
	})
	if err != nil {
		return models.Pedido{}, err
	}

	evento := messaging.PedidoConfirmado{
		Tipo:       "pedido.confirmado",
		PedidoID:   pedido.ID,
		ClienteID:  pedido.ClienteID,
		ProductoID: pedido.ProductoID,
	}

	if err := s.Publisher.PublicarPedidoConfirmado(evento); err != nil {
		return models.Pedido{}, err
	}

	return pedido, nil
}
