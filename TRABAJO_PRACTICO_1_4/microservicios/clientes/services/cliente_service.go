package services

import (
	"errors"

	"clientes/models"
	"clientes/repositories"
)

type ClienteService struct {
	Repo repositories.ClientesRepo
}

func (s *ClienteService) CrearCliente(cliente models.Cliente) (models.Cliente, error) {
	if cliente.ID == "" || cliente.Nombre == "" {
		return models.Cliente{}, errors.New("id y nombre son obligatorios")
	}

	if err := s.Repo.CreateCliente(cliente); err != nil {
		return models.Cliente{}, err
	}
	return cliente, nil
}

func (s *ClienteService) ObtenerCliente(id string) (models.Cliente, error) {
	if id == "" {
		return models.Cliente{}, errors.New("el ID no puede estar vacío")
	}
	return s.Repo.GetClienteByID(id)
}
