package app

import (
	"github.com/gin-gonic/gin"
	http2 "github.com/nuzurie/shopify/inventory/delivery/http"
	"github.com/nuzurie/shopify/item/delivery/http"
)

func mapItemUrls(handler *http.ItemHandler, r *gin.Engine) {
	r.GET("/items", handler.GetAll)
	r.POST("/items", handler.Create)
	r.PUT("/items/:id", handler.Update)
	r.DELETE("/items/:id", handler.Delete)
}

func mapInventoryUrls(handler *http2.InventoryHandler, r *gin.Engine) {
	r.GET("/inventory", handler.GetAll)
	r.GET("/inventory/:id", handler.GetInventoryForItem)
	r.POST("/inventory", handler.CreateOrUpdate)
	r.DELETE("/inventory/:id", handler.Delete)
}
