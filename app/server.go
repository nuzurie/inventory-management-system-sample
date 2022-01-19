package app

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	http2 "github.com/nuzurie/shopify/inventory/delivery/http"
	repository2 "github.com/nuzurie/shopify/inventory/repository"
	usecase2 "github.com/nuzurie/shopify/inventory/usecase"
	"github.com/nuzurie/shopify/item/delivery/http"
	"github.com/nuzurie/shopify/item/repository"
	"github.com/nuzurie/shopify/item/usecase"
	"log"
	"os"
	"time"
)

func Server(itemHandler *http.ItemHandler, inventoryHandler *http2.InventoryHandler) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	mapItemUrls(itemHandler, router)
	mapInventoryUrls(inventoryHandler, router)
	return router
}

func Start() {
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(os.Getenv("DATABASE_URL"))
		log.Fatalln("db failed", err)
	}

	itemRepository, err := repository.NewItemRepository(pool)
	if err != nil {
		log.Fatalln("Failed to initialize item table ", err)
	}
	itemUseCase := usecase.NewItemUseCase(itemRepository, time.Second)
	itemHandler := http.NewItemHandler(itemUseCase)

	inventoryRepository, err := repository2.NewInventoryRepository(pool)
	if err != nil {
		log.Fatalln("Failed to initialize item table ", err)
	}
	inventoryUseCase := usecase2.NewInventoryUseCase(itemRepository, inventoryRepository, time.Second*300)
	inventoryHandler := http2.NewInventoryHandler(inventoryUseCase)

	router := Server(itemHandler, inventoryHandler)
	router.Run()
}
