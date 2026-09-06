package repositories

import "pedidos/models"

type ProductosRepo interface {
	ListarProductos() ([]models.Producto, error)
}

// ProductosMemRepo simula la base de datos de productos en memoria.
type ProductosMemRepo struct {
	productos []models.Producto
}

func NewProductosMemRepo() *ProductosMemRepo {
	return &ProductosMemRepo{
		productos: []models.Producto{
			{ID: "P-1", Nombre: "Auriculares", Stock: 10},
			{ID: "P-2", Nombre: "Teclado", Stock: 8},
		},
	}
}

func (r *ProductosMemRepo) ListarProductos() ([]models.Producto, error) {
	return r.productos, nil
}
