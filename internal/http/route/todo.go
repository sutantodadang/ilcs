package route

import (
	"ilcs/internal/app/todo"

	"github.com/gin-gonic/gin"
)

func RegisterTodoRoute(app *gin.Engine, handler todo.ITodoHandler) {
	todoRoute := app.Group("/api/v1")
	todoRoute.POST("/tasks", handler.CreateTodo)
	todoRoute.GET("/tasks", handler.ListTodo)
	todoRoute.GET("/tasks/:id", handler.GetTodoById)
	todoRoute.PUT("/tasks/:id", handler.UpdateTodo)
	todoRoute.DELETE("/tasks/:id", handler.DeleteTodo)

}
