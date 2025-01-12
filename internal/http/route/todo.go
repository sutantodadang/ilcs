package route

import (
	"ilcs/internal/app/todo"
	"ilcs/internal/http/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterTodoRoute(app *gin.Engine, handler todo.ITodoHandler) {
	todoRoute := app.Group("/api/v1")
	todoRoute.POST("/tasks", middlewares.Auth(), handler.CreateTodo)
	todoRoute.GET("/tasks", middlewares.Auth(), handler.ListTodo)
	todoRoute.GET("/tasks/:id", middlewares.Auth(), handler.GetTodoById)
	todoRoute.PUT("/tasks/:id", middlewares.Auth(), handler.UpdateTodo)
	todoRoute.DELETE("/tasks/:id", middlewares.Auth(), handler.DeleteTodo)
	todoRoute.GET("/token", handler.GetToken)

}
