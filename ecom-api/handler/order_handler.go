package handler

import (
	"net/http"
	"strconv"

	"github.com/codepnw/microservice-ecommerce/ecom-api/store"
	"github.com/codepnw/microservice-ecommerce/token"
	"github.com/gin-gonic/gin"
)

func (h *handler) createOrder(c *gin.Context) {
	var o OrderReq
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get Context
	claims, exists := c.Get(claimsKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	so := toStoreOrder(o)
	so.UserID = claims.(*token.UserClaims).ID

	created, err := h.server.CreateOrder(c.Request.Context(), so)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := toOrderRes(created)
	c.JSON(http.StatusCreated, res)
}

func (h *handler) getOrder(c *gin.Context) {
	// Get Context
	claims, exists := c.Get(claimsKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	id := claims.(*token.UserClaims).ID

	order, err := h.server.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := toOrderRes(order)
	c.JSON(http.StatusOK, res)
}

func (h *handler) listOrders(c *gin.Context) {
	orders, err := h.server.ListOrder(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var res []OrderRes
	for _, o := range orders {
		res = append(res, toOrderRes(&o))
	}

	c.JSON(http.StatusOK, res)
}

func (h *handler) deleteOrder(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error pasing ID"})
		return
	}

	if err := h.server.DeleteOrder(c.Request.Context(), idInt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func toStoreOrder(o OrderReq) *store.Order {
	return &store.Order{
		PaymentMethod: o.PaymentMethod,
		TaxPrice:      o.TaxPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		Items:         toStoreOrderItem(o.Items),
	}
}

func toStoreOrderItem(items []*OrderItem) []store.OrderItem {
	var res []store.OrderItem

	for _, i := range items {
		res = append(res, store.OrderItem{
			Name:      i.Name,
			Quantity:  i.Quantity,
			Image:     i.Image,
			Price:     i.Price,
			ProductID: i.ProductID,
		})
	}
	return res
}

func toOrderRes(o *store.Order) OrderRes {
	return OrderRes{
		ID:            o.ID,
		Items:         toOrderItems(o.Items),
		PaymentMethod: o.PaymentMethod,
		TaxPrice:      o.TaxPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		CreatedAt:     o.CreatedAt,
		UpdatedAt:     o.UpdatedAt,
	}
}

func toOrderItems(items []store.OrderItem) []OrderItem {
	var res []OrderItem
	for _, i := range items {
		res = append(res, OrderItem{
			Name:      i.Name,
			Quantity:  i.Quantity,
			Image:     i.Image,
			Price:     i.Price,
			ProductID: i.ProductID,
		})
	}
	return res
}
