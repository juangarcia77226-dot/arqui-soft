package main

import (
	"context"
	"log"
	"time"

	// Importamos el router (usamos Gin como ejemplo)
	"github.com/gin-gonic/gin"

	// Cache local (in-memory) para el patrón Cache-Aside
	"github.com/karlseguin/ccache/v2"

	// Importamos nuestros paquetes internos
	"main/controllers/items"
	"main/database"
	repo "main/repositories/items"
	services "main/services/items"
)

func main() {
	// 1. Levantamos la conexión a la base de datos
	// Comentamos esta linea para probar test
	client := database.ConnectDB()

	// 2. Aseguramos el cierre de la conexión al detener la API
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal("Error al desconectar MongoDB:", err)
		}
	}()

	log.Println("Conexión exitosa. Configurando capas de la aplicación...")

	// ---------------------------------------------------------
	// 3. INYECCIÓN DE DEPENDENCIAS (El Cableado)
	// ---------------------------------------------------------

	// // A. Repositorio: Le pasamos el cliente real de MongoDB
	mongoRepo := repo.ItemsMongoDB{
		Client: client,
	}

	// A.1 Cache-Aside: envolvemos el repositorio real con una caché local
	// (ccache). El Service sigue viendo un simple ItemsRepo, no sabe que
	// hay una caché en el medio.
	itemsCache := ccache.New(ccache.Configure())
	defer itemsCache.Stop()

	itemsRepo := repo.ItemsCachedRepo{
		Repo:  mongoRepo,
		Cache: itemsCache,
		TTL:   60 * time.Second,
	}

	// B. Servicio: Le inyectamos el repositorio (cumple la interfaz ItemsRepo)
	itemsService := services.ItemsService{
		Repo: itemsRepo,
	}

	// 2. Preparamos el mapa vacío para la RAM
	// dbEnMemoria := map[string]modelItems.ItemModel{
	// 	"1": {
	// 		ID:    "1",
	// 		Title: "Laptop",
	// 		Price: 1500,
	// 	},
	// }
	// // 3. Instanciamos el repo Mock
	// itemsRepo := repo.ItemsMock{
	// 	MockDB: dbEnMemoria,
	// }

	// // 4. Se lo pasamos al Service. ¡Magia! El Service lo acepta igual.
	// itemsService := services.ItemsService{
	// 	Repo: itemsRepo,
	// }

	// C. Controlador: Le inyectamos el servicio con la lógica de negocio
	itemsController := items.ItemsController{
		Service: itemsService,
	}

	// ---------------------------------------------------------
	// 4. RUTEO Y SERVIDOR
	// ---------------------------------------------------------

	// Inicializamos el framework web Gin
	router := gin.Default()

	// Definimos los endpoints y a qué función del controlador apuntan
	router.GET("/items/:id", itemsController.GetByID)

	// router.POST("/items", itemsController.Create)
	// router.PUT("/items/:id", itemsController.Update)

	// Levantamos el servidor en el puerto 8080
	log.Println("Servidor corriendo en http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error al arrancar el servidor:", err)
	}
}

// [ Cliente / Frontend ]
//         │  ▲
//     DTO │  │ DTO
//         ▼  │
// [ Capa de Controladores ] <─── (Mapeo: DTO 🔄 Model)
//         │  ▲
//   Model │  │ Model
//         ▼  │
// [ Capa de Negocio / Service ]
//         │  ▲
//   Model │  │ Model
//         ▼  │
// [ Capa de Repositorio (Repository) ]
//         │  ▲
//         │  │
//         ▼  │
// [ Capa de Acceso a Datos (Model) ] <─── (Traduce Model a SQL/Query)
//         │  ▲
//         │  │ Registros puros (Rows)
//         ▼  │
// [ Base de Datos ]
