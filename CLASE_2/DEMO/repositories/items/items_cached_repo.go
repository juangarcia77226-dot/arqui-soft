package items

import (
	"log"
	"main/models/items"
	"time"

	"github.com/karlseguin/ccache/v2"
)

// ItemsCachedRepo implementa el patrón Cache-Aside: envuelve a otro
// ItemsRepo (por ejemplo ItemsMongoDB) y resuelve las lecturas contra una
// caché local (ccache) antes de golpear la base de datos.
//
//	Cache hit  -> devolvemos el dato desde memoria
//	Cache miss -> consultamos la BD real y guardamos el resultado en caché
//
// Como implementa la misma interfaz ItemsRepo, el Service y el Controller
// no se enteran de que están hablando con una caché: sólo reciben "un
// repositorio" más.
type ItemsCachedRepo struct {
	Repo  ItemsRepo     // repositorio real al que se le hace el cache-aside
	Cache *ccache.Cache // caché local en memoria (ccache)
	TTL   time.Duration // tiempo de expiración de cada entrada
}

func (r ItemsCachedRepo) GetItemByID(id string) (items.ItemModel, error) {
	start := time.Now()

	// 1. Cache hit: si el item está en caché y no expiró, lo devolvemos
	if value := r.Cache.Get(id); value != nil && !value.Expired() {
		log.Printf("[cache] HIT  id=%s (%s)", id, time.Since(start))
		return value.Value().(items.ItemModel), nil
	}

	// 2. Cache miss: consultamos la base de datos
	item, err := r.Repo.GetItemByID(id)
	if err != nil {
		return items.ItemModel{}, err
	}

	// 3. Guardamos el resultado en caché para las próximas consultas
	r.Cache.Set(id, item, r.TTL)
	log.Printf("[cache] MISS id=%s (%s)", id, time.Since(start))

	return item, nil
}

func (r ItemsCachedRepo) CreateItem(item items.ItemModel) error {
	// Al escribir invalidamos la entrada en caché para no servir datos
	// desactualizados en la próxima lectura.
	r.Cache.Delete(item.ID)
	return r.Repo.CreateItem(item)
}
