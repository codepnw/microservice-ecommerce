package handler

import "github.com/gin-gonic/gin"

var r *gin.Engine

func RegisterRoutes(handler *handler) *gin.Engine {
	r = gin.Default()

	products := r.Group("/products")
	{
		products.POST("/", handler.createProduct)
		products.GET("/", handler.listProducts)

		productID := products.Group("/:id")
		{
			productID.GET("", handler.getProduct)
			productID.PATCH("", handler.updateProduct)
			productID.DELETE("", handler.deleteProduct)
		}
	}

	orders := r.Group("/orders")
	{
		orders.POST("/", handler.createOrder)
		orders.GET("/", handler.listOrders)

		orderID := orders.Group("/:id")
		{
			orderID.GET("", handler.getOrder)
			orderID.DELETE("", handler.deleteOrder)
		}
	}

	users := r.Group("/users")
	{
		users.POST("/", handler.createUser)
		users.GET("/", handler.listUsers)
		users.PATCH("/", handler.updateUser)
		users.DELETE("/:id", handler.deleteUser)
	}

	// Auth
	r.POST("/login", handler.loginUser)
	r.POST("/logout", handler.logoutUser)

	// Token
	r.POST("/token/renew", handler.renewAccessToken)
	r.POST("/token/revoke/:id", handler.revokeSession)

	return r
}

func Start(addr string) *gin.Engine {
	r.Run(":" + addr)
	return r
}
