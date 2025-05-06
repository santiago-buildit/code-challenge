package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/santiago-buildit/code-challenge/backend/internal/handlers"
)

func RegisterBookRoutes(router *gin.Engine, handler *handlers.BookHandler) {
	group := router.Group("/books")
	{
		// CRUD operations
		group.POST("", handler.CreateBook)
		group.POST("/list", handler.ListBooks) // POST to support pagination and filtering
		group.GET("/:id", handler.GetBook)
		group.PUT("/:id", handler.UpdateBook)
		group.DELETE("/:id", handler.DeleteBook)

		// Status operations
		group.PUT("/:id/checkout", handler.CheckoutBook)
		group.PUT("/:id/checkin", handler.CheckinBook)

		// History
		group.GET("/:id/details", handler.GetBookWithHistory)
	}
}
