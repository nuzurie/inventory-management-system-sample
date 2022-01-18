package app

import (
	"github.com/gin-gonic/gin"
	"github.com/nuzurie/shopify/item/delivery/http"
)

func mapItemUrls(handler *http.ItemHandler, r *gin.Engine) {
	r.GET("/items", handler.GetAll)
	r.POST("/items", handler.Create)
	r.PUT("/items/:id", handler.Update)
	r.DELETE("/items/:id", handler.Delete)
}
