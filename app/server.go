package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nuzurie/shopify/item/delivery/http"
	"github.com/nuzurie/shopify/item/repository"
	"github.com/nuzurie/shopify/item/usecase"
	"log"
	"os"
	"time"
)

func Server(handler *http.ItemHandler) *gin.Engine {
	router := gin.Default()
	mapItemUrls(handler, router)
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

	router := Server(itemHandler)
	router.Run()
}
