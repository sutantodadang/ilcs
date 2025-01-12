package todo

type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	DueDate     string `json:"due_date" binding:"required,datetime=2006-01-02"`
}

type Todo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDate     string `json:"due_date"`
}

type ListTodoRequestParams struct {
	Page   *int    `form:"page"`
	Limit  *int    `form:"limit"`
	Status *string `form:"status"`
	Search *string `form:"search"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"required,oneof=pending,completed"`
	DueDate     string `json:"due_date" binding:"required,datetime=2006-01-02"`
}
