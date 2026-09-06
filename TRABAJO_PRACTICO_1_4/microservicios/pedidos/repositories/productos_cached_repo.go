package repositories

import (
	"log"
	"sync"
	"time"

	"pedidos/models"
)

// ProductosCachedRepo implementa el patrón Cache-Aside sobre otro
// ProductosRepo: sirve el listado de productos desde memoria mientras no
// haya expirado el TTL, y sólo vuelve a consultar el repositorio real
// (Cache Miss) cuando la caché está vacía o vencida.
type ProductosCachedRepo struct {
	Repo ProductosRepo
	TTL  time.Duration

	mu        sync.Mutex
	productos []models.Producto
	expiresAt time.Time
}

func NewProductosCachedRepo(repo ProductosRepo, ttl time.Duration) *ProductosCachedRepo {
	return &ProductosCachedRepo{Repo: repo, TTL: ttl}
}

func (r *ProductosCachedRepo) ListarProductos() ([]models.Producto, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Cache hit: todavía no expiró
	if r.productos != nil && time.Now().Before(r.expiresAt) {
		log.Println("[cache] HIT listado de productos")
		return r.productos, nil
	}

	// Cache miss: consultamos el repositorio real
	productos, err := r.Repo.ListarProductos()
	if err != nil {
		return nil, err
	}

	r.productos = productos
	r.expiresAt = time.Now().Add(r.TTL)
	log.Println("[cache] MISS listado de productos, recargado desde el repositorio")

	return productos, nil
}
