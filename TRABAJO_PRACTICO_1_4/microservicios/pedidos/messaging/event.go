package messaging

// PedidoConfirmado es el evento que pedidos publica cuando confirma una
// compra. En un caso real, logística lo consumiría para preparar el envío.
type PedidoConfirmado struct {
	Tipo       string `json:"tipo"`
	PedidoID   string `json:"pedido_id"`
	ClienteID  string `json:"cliente_id"`
	ProductoID string `json:"producto_id"`
}
