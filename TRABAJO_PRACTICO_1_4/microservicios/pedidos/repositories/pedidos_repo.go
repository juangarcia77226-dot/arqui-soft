package repositories

import (
	"fmt"
	"sync"

	"pedidos/models"
)

type PedidosRepo interface {
	CrearPedido(pedido models.Pedido) (models.Pedido, error)
}

// PedidosMemRepo simula la base de datos de pedidos en memoria.
type PedidosMemRepo struct {
	mu       sync.Mutex
	pedidos  map[string]models.Pedido
	contador int
}

func NewPedidosMemRepo() *PedidosMemRepo {
	return &PedidosMemRepo{pedidos: make(map[string]models.Pedido)}
}

func (r *PedidosMemRepo) CrearPedido(pedido models.Pedido) (models.Pedido, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.contador++
	pedido.ID = fmt.Sprintf("PED-%d", r.contador)
	r.pedidos[pedido.ID] = pedido

	return pedido, nil
}
