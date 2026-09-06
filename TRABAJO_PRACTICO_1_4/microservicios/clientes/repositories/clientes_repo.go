package repositories

import (
	"errors"
	"sync"

	"clientes/models"
)

type ClientesRepo interface {
	CreateCliente(cliente models.Cliente) error
	GetClienteByID(id string) (models.Cliente, error)
}

// ClientesMemRepo simula la base de datos de clientes en memoria.
type ClientesMemRepo struct {
	mu   sync.RWMutex
	data map[string]models.Cliente
}

func NewClientesMemRepo() *ClientesMemRepo {
	return &ClientesMemRepo{data: make(map[string]models.Cliente)}
}

func (r *ClientesMemRepo) CreateCliente(cliente models.Cliente) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[cliente.ID] = cliente
	return nil
}

func (r *ClientesMemRepo) GetClienteByID(id string) (models.Cliente, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cliente, existe := r.data[id]
	if !existe {
		return models.Cliente{}, errors.New("cliente no encontrado")
	}
	return cliente, nil
}
