package models

type Pedido struct {
	ID         string `json:"id"`
	ClienteID  string `json:"cliente_id"`
	ProductoID string `json:"producto_id"`
}
