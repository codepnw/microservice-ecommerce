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

		users.POST("/login", handler.loginUser)
		users.POST("/logout", handler.logoutUser)

		token := users.Group("/token")
		{
			token.POST("/renew", handler.renewAccessToken)
			token.POST("/revoke/:id", handler.revokeSession)
		}
	}

	return r
}

func Start(addr string) *gin.Engine {
	r.Run(":" + addr)
	return r
}
