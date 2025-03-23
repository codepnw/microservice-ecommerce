package handler

import "github.com/gin-gonic/gin"

var r *gin.Engine

func RegisterRoutes(handler *handler) *gin.Engine {
	r = gin.Default()
	tokenMaker := handler.TokenMaker

	products := r.Group("/products")
	{
		products.POST("/", GetAdminMiddlewareFunc(tokenMaker), handler.createProduct)
		products.GET("/", handler.listProducts)

		productID := products.Group("/:id")
		{
			productID.GET("", handler.getProduct)
			productID.PATCH("", GetAdminMiddlewareFunc(tokenMaker), handler.updateProduct)
			productID.DELETE("", GetAdminMiddlewareFunc(tokenMaker), handler.deleteProduct)
		}
	}

	orders := r.Group("/orders")
	{
		orders.Use(GetAuthMiddlewareFunc(tokenMaker))

		orders.GET("/", GetAdminMiddlewareFunc(tokenMaker), handler.listOrders)
		orders.POST("/", handler.createOrder)
		orders.GET("/myorder", handler.getOrder)
		orders.DELETE("/:id", handler.deleteOrder)
	}

	users := r.Group("/users")
	{
		users.POST("/", handler.createUser)

		users.GET("/", GetAdminMiddlewareFunc(tokenMaker), handler.listUsers)
		users.DELETE("/:id", GetAdminMiddlewareFunc(tokenMaker), handler.deleteUser)

		users.PATCH("/", GetAuthMiddlewareFunc(tokenMaker), handler.updateUser)
	}

	// Auth
	r.POST("/login", handler.loginUser)
	r.POST("/logout", GetAuthMiddlewareFunc(tokenMaker), handler.logoutUser)

	// Token
	r.POST("/token/renew", GetAuthMiddlewareFunc(tokenMaker), handler.renewAccessToken)
	r.POST("/token/revoke", GetAuthMiddlewareFunc(tokenMaker), handler.revokeSession)

	return r
}

func Start(addr string) *gin.Engine {
	r.Run(":" + addr)
	return r
}
