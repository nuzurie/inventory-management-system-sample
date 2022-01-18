package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nuzurie/shopify/domain"
	"github.com/nuzurie/shopify/specification"
	"github.com/nuzurie/shopify/utils/errors"
	"net/http"
	"reflect"
	"strconv"
)

type ItemHandler struct {
	useCase domain.ItemUseCase
}

func NewItemHandler(useCase domain.ItemUseCase) *ItemHandler {
	return &ItemHandler{useCase: useCase}
}

func (h *ItemHandler) GetAll(c *gin.Context) {
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

	spec := specification.NewItemSpecification(name, description, minPrice, maxPrice)

	ctx := c.Request.Context()
	items, err := h.useCase.GetAll(ctx, int(count), int(offset), spec)
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

func (h *ItemHandler) Create(c *gin.Context) {
	var item domain.Item
	err := c.ShouldBind(&item)
	if err != nil || reflect.DeepEqual(&domain.Item{}, item) {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid review body"))
		return
	}

	ctx := c.Request.Context()
	createdItem, err := h.useCase.Create(ctx, &item)
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

	c.JSON(http.StatusCreated, createdItem)
}

func (h *ItemHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("id not provided"))
		return
	}

	var item domain.Item
	err := c.ShouldBind(&item)
	if err != nil || reflect.DeepEqual(&domain.Item{}, item) {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid review body"))
		return
	}

	item.ID = id
	ctx := c.Request.Context()
	updated, err := h.useCase.Update(ctx, &item)
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

	c.JSON(http.StatusOK, updated)
}

func (h *ItemHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("id not provided"))
		return
	}

	ctx := c.Request.Context()
	err := h.useCase.Delete(ctx, id)
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

	c.JSON(http.StatusNoContent, gin.H{"message": "deleted successfully"})
}
