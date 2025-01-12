package todo

import (
	"context"
	"testing"
	"time"

	"ilcs/internal/app/todo"
	"ilcs/internal/repositories"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) InsertTodo(ctx context.Context, params repositories.InsertTodoParams) (repositories.Todo, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(repositories.Todo), args.Error(1)
}

func (m *MockRepo) ListTodo(ctx context.Context, params repositories.ListTodoParams) ([]repositories.ListTodoRow, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]repositories.ListTodoRow), args.Error(1)
}

func (m *MockRepo) CountTodo(ctx context.Context, params repositories.CountTodoParams) (int64, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepo) GetTodoById(ctx context.Context, id pgtype.UUID) (repositories.GetTodoByIdRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(repositories.GetTodoByIdRow), args.Error(1)
}

func (m *MockRepo) UpdateTodo(ctx context.Context, params repositories.UpdateTodoParams) (repositories.Todo, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(repositories.Todo), args.Error(1)
}

func (m *MockRepo) DeleteTodo(ctx context.Context, id pgtype.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return redis.NewStringResult(args.String(0), args.Error(1))
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return redis.NewStatusResult(args.String(0), args.Error(1))
}

func TestCreateTodo_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	mockRedisClient := new(MockRedisClient)
	service := todo.NewTodoService(mockRepo, mockRedisClient)

	req := todo.CreateTodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
		DueDate:     "2025-01-01",
	}

	expectedTodo := repositories.Todo{
		ID:          pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Title:       req.Title,
		Description: pgtype.Text{String: req.Description, Valid: true},
		DueDate:     pgtype.Date{Time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true},
	}

	mockRepo.On("InsertTodo", mock.Anything, mock.Anything).Return(expectedTodo, nil)

	todo, err := service.CreateTodo(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, expectedTodo, todo)
	mockRepo.AssertExpectations(t)
}

func TestCreateTodo_InvalidDate(t *testing.T) {
	mockRepo := new(MockRepo)
	mockRedisClient := new(MockRedisClient)
	service := todo.NewTodoService(mockRepo, mockRedisClient)

	req := todo.CreateTodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
		DueDate:     "invalid-date",
	}

	_, err := service.CreateTodo(context.Background(), req)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "InsertTodo")
}

func TestGetListTodos_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	mockRedisClient := new(MockRedisClient)
	service := todo.NewTodoService(mockRepo, mockRedisClient)

	req := todo.ListTodoRequestParams{
		Page:  func(i int) *int { return &i }(1),
		Limit: func(i int) *int { return &i }(10),
	}

	returnTodos := []repositories.ListTodoRow{
		{
			ID:          pgtype.UUID{Bytes: uuid.New(), Valid: true},
			Title:       "Test Todo 1",
			Description: pgtype.Text{String: "Test Description 1", Valid: true},
			DueDate:     pgtype.Date{Time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true},
		},
		{
			ID:          pgtype.UUID{Bytes: uuid.New(), Valid: true},
			Title:       "Test Todo 2",
			Description: pgtype.Text{String: "Test Description 2", Valid: true},
			DueDate:     pgtype.Date{Time: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC), Valid: true},
		},
	}

	expectedTodos := []todo.Todo{
		{
			ID:          returnTodos[0].ID.String(),
			Title:       returnTodos[0].Title,
			Description: returnTodos[0].Description.String,
			DueDate:     returnTodos[0].DueDate.Time.Format("2006-01-02"),
		},
		{
			ID:          returnTodos[1].ID.String(),
			Title:       returnTodos[1].Title,
			Description: returnTodos[1].Description.String,
			DueDate:     returnTodos[1].DueDate.Time.Format("2006-01-02"),
		},
	}

	mockRepo.On("ListTodo", mock.Anything, mock.Anything).Return(returnTodos, nil)
	mockRepo.On("CountTodo", mock.Anything, mock.Anything).Return(int64(len(returnTodos)), nil)

	todos, count, page, limit, err := service.GetListTodos(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, expectedTodos, todos)
	assert.Equal(t, int64(len(expectedTodos)), count)
	assert.Equal(t, 1, page)
	assert.Equal(t, 10, limit)
	mockRepo.AssertExpectations(t)
}

func TestGetTodo_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	mockRedisClient := new(MockRedisClient)
	service := todo.NewTodoService(mockRepo, mockRedisClient)

	id := uuid.New().String()

	returnTodo := repositories.GetTodoByIdRow{
		ID:          pgtype.UUID{Bytes: uuid.MustParse(id), Valid: true},
		Title:       "Test Todo",
		Description: pgtype.Text{String: "Test Description", Valid: true},
		DueDate:     pgtype.Date{Time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true},
	}

	expectedTodo := todo.Todo{

		ID:          returnTodo.ID.String(),
		Title:       returnTodo.Title,
		Description: returnTodo.Description.String,
		DueDate:     returnTodo.DueDate.Time.Format("2006-01-02"),
	}

	mockRepo.On("GetTodoById", mock.Anything, mock.Anything).Return(returnTodo, nil)
	mockRedisClient.On("Get", mock.Anything, "todo:"+id).Return("", redis.Nil)
	mockRedisClient.On("Set", mock.Anything, "todo:"+id, mock.Anything, mock.Anything).Return("OK", nil)

	todo, err := service.GetTodo(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, expectedTodo, todo)
	mockRepo.AssertExpectations(t)
	mockRedisClient.AssertExpectations(t)
}

func TestUpdateTodo_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	mockRedisClient := new(MockRedisClient)
	service := todo.NewTodoService(mockRepo, mockRedisClient)

	id := uuid.New().String()

	req := todo.UpdateTodoRequest{
		Title:       "Updated Todo",
		Description: "Updated Description",
		DueDate:     "2025-01-02",
	}

	expectedTodo := repositories.Todo{
		ID:          pgtype.UUID{Bytes: uuid.MustParse(id), Valid: true},
		Title:       req.Title,
		Description: pgtype.Text{String: req.Description, Valid: true},
		DueDate:     pgtype.Date{Time: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC), Valid: true},
	}

	mockRepo.On("UpdateTodo", mock.Anything, mock.Anything).Return(expectedTodo, nil)

	todo, err := service.UpdateTodo(context.Background(), req, id)

	assert.NoError(t, err)
	assert.Equal(t, expectedTodo, todo)
	mockRepo.AssertExpectations(t)
}

func TestDeleteTodo_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	mockRedisClient := new(MockRedisClient)
	service := todo.NewTodoService(mockRepo, mockRedisClient)

	id := uuid.New().String()

	mockRepo.On("DeleteTodo", mock.Anything, mock.Anything).Return(nil)

	err := service.DeleteTodo(context.Background(), id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
