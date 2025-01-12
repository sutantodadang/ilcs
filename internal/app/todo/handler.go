package todo

import (
	"ilcs/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ITodoHandler interface {
	CreateTodo(c *gin.Context)
	ListTodo(c *gin.Context)
	GetTodoById(c *gin.Context)
	UpdateTodo(c *gin.Context)
	DeleteTodo(c *gin.Context)
	GetToken(c *gin.Context)
}

type TodoHandler struct {
	service ITodoService
}

func NewTodoHandler(service ITodoService) *TodoHandler {
	return &TodoHandler{
		service: service,
	}
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {

	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(400, gin.H{"error": utils.NewValidationError(errs)})
			return
		}

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.service.CreateTodo(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Task created successfully", "task": todo})
}

func (h *TodoHandler) ListTodo(c *gin.Context) {

	var req ListTodoRequestParams
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	todos, countData, currentPage, currentLimit, err := h.service.GetListTodos(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var response struct {
		Tasks      []Todo `json:"tasks"`
		Pagination struct {
			CurrentPage int   `json:"current_page"`
			TotalPage   int   `json:"total_page"`
			TotalTasks  int64 `json:"total_tasks"`
		} `json:"pagination"`
	}

	response.Tasks = todos
	response.Pagination.CurrentPage = currentPage
	response.Pagination.TotalPage = int(countData) / currentLimit
	response.Pagination.TotalTasks = countData

	c.JSON(200, response)
}

func (h *TodoHandler) GetTodoById(c *gin.Context) {

	id := c.Param("id")

	if err := utils.ValidateId(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return

	}

	todo, err := h.service.GetTodo(c, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, todo)

}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {

	id := c.Param("id")

	if err := utils.ValidateId(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return

	}

	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(400, gin.H{"error": utils.NewValidationError(errs)})
			return
		}

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.service.UpdateTodo(c, req, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Task updated successfully", "task": todo})
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {

	id := c.Param("id")

	if err := utils.ValidateId(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return

	}

	err := h.service.DeleteTodo(c, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Task deleted successfully"})
}

func (h *TodoHandler) GetToken(c *gin.Context) {

	token, err := h.service.GetToken(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token})
}
