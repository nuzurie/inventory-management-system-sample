package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nuzurie/shopify/domain"
	"github.com/nuzurie/shopify/specification"
	"github.com/nuzurie/shopify/utils/errors"
	"net/http"
	"reflect"
	"strconv"
)

type InventoryHandler struct {
	useCase domain.InventoryUseCase
}

func NewInventoryHandler(useCase domain.InventoryUseCase) *InventoryHandler {
	return &InventoryHandler{useCase: useCase}
}

func (h *InventoryHandler) GetAll(c *gin.Context) {
	name, _ := c.GetQuery("name")
	description, _ := c.GetQuery("description")
	minPriceQuery, _ := c.GetQuery("min-price")
	minPrice, err := strconv.ParseFloat(minPriceQuery, 64)
	if err != nil {
		minPrice = 0
	}
	maxPriceQuery, ok := c.GetQuery("max-price")
	var maxPrice float64
	if ok {
		maxPrice, err = strconv.ParseFloat(maxPriceQuery, 64)
		if err != nil {
			maxPrice = -1
		}
	} else {
		maxPrice = -1
	}

	itemSpec := specification.NewItemSpecification(name, description, minPrice, maxPrice)

	minQuantityQuery, _ := c.GetQuery("min-quantity")
	minQuantity, err := strconv.ParseInt(minQuantityQuery, 10, 64)
	if err != nil {
		minQuantity = 0
	}
	maxQuantityQuery, ok := c.GetQuery("max-quantity")
	var maxQuantity int64
	if ok {
		maxQuantity, err = strconv.ParseInt(maxQuantityQuery, 10, 64)
		if err != nil {
			maxQuantity = -1
		}
	} else {
		maxQuantity = -1
	}

	inventorySpec := specification.NewInventorySpecification(int(minQuantity), int(maxQuantity), itemSpec)

	var count int64
	if countQuery, ok := c.GetQuery("count"); ok {
		count, _ = strconv.ParseInt(countQuery, 10, 64)
	} else {
		count = 20
	}

	var offset int64
	if offsetQuery, ok := c.GetQuery("offset"); ok {
		offset, _ = strconv.ParseInt(offsetQuery, 10, 64)
	}

	ctx := c.Request.Context()
	items, err := h.useCase.GetAll(ctx, int(count), int(offset), inventorySpec)
	if err != nil {
		switch v := err.(type) {
		case *errors.RestError:
			c.JSON(v.Code, v)
			return
		default:
			c.JSON(http.StatusInternalServerError, errors.NewInternalServerError(err.Error()))
			return
		}
	}

	c.JSON(http.StatusOK, items)
}

func (h *InventoryHandler) GetInventoryForItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("id not provided"))
		return
	}

	ctx := c.Request.Context()
	inventory, err := h.useCase.GetInventoryForItem(ctx, id)
	if err != nil {
		switch v := err.(type) {
		case *errors.RestError:
			c.JSON(v.Code, v)
			return
		default:
			c.JSON(http.StatusInternalServerError, errors.NewInternalServerError(err.Error()))
			return
		}
	}

	c.JSON(http.StatusNoContent, inventory)
}

func (h *InventoryHandler) CreateOrUpdate(c *gin.Context) {
	var inventory domain.InventoryItem
	err := c.ShouldBind(&inventory)
	if err != nil || reflect.DeepEqual(&domain.InventoryItem{}, inventory) {
		fmt.Println(inventory)
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid inventory body"))
		return
	}

	ctx := c.Request.Context()
	createdInventory, err := h.useCase.UpdateInventoryItem(ctx, &inventory)
	if err != nil {
		switch v := err.(type) {
		case *errors.RestError:
			c.JSON(v.Code, v)
			return
		default:
			c.JSON(http.StatusInternalServerError, errors.NewInternalServerError(err.Error()))
			return
		}
	}

	c.JSON(http.StatusCreated, createdInventory)
}

func (h *InventoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("id not provided"))
		return
	}

	ctx := c.Request.Context()
	err := h.useCase.DeleteItem(ctx, id)
	if err != nil {
		switch v := err.(type) {
		case *errors.RestError:
			c.JSON(v.Code, v)
			return
		default:
			c.JSON(http.StatusInternalServerError, errors.NewInternalServerError(err.Error()))
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}
