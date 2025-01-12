package todo

import (
	"context"
	"ilcs/internal/repositories"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

type ITodoService interface {
	CreateTodo(ctx context.Context, req CreateTodoRequest) (todo repositories.Todo, err error)
	GetListTodos(ctx context.Context, req ListTodoRequestParams) (todos []Todo, countData int64, page, limit int, err error)
	GetTodo(ctx context.Context, id string) (todo Todo, err error)
	UpdateTodo(ctx context.Context, req UpdateTodoRequest, id string) (todo repositories.Todo, err error)
	DeleteTodo(ctx context.Context, id string) (err error)
}

type TodoService struct {
	repo repositories.Querier
}

func NewTodoService(repo repositories.Querier) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

func (s *TodoService) CreateTodo(ctx context.Context, req CreateTodoRequest) (todo repositories.Todo, err error) {

	var wg sync.WaitGroup

	errChan := make(chan error, 1)
	todoChan := make(chan repositories.Todo, 1)

	wg.Add(1)
	go func(ctx context.Context, req CreateTodoRequest) {
		defer wg.Done()

		id, err := uuid.NewV7()
		if err != nil {
			log.Error().Err(err).Send()
			errChan <- err
			return
		}

		timeDate, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			log.Error().Err(err).Send()
			errChan <- err
			return
		}

		todo, err = s.repo.InsertTodo(ctx, repositories.InsertTodoParams{
			ID:          pgtype.UUID{Bytes: id, Valid: true},
			Title:       req.Title,
			Description: pgtype.Text{String: req.Description, Valid: true},
			DueDate:     pgtype.Date{Time: timeDate, Valid: true},
		})

		if err != nil {
			log.Error().Err(err).Send()
			errChan <- err
			return
		}

		todoChan <- todo

	}(ctx, req)

	wg.Wait()

	select {
	case err = <-errChan:
		if err != nil {
			log.Error().Err(err).Send()
			return
		}
	case todo = <-todoChan:
		return todo, nil
	case <-ctx.Done():
		err = ctx.Err()
		log.Error().Err(err).Send()
		return
	}

	close(errChan)
	close(todoChan)

	return
}

func (s *TodoService) GetListTodos(ctx context.Context, req ListTodoRequestParams) (todos []Todo, countData int64, page, limit int, err error) {

	if req.Page == nil {
		req.Page = new(int)
		*req.Page = 1
		page = *req.Page

	} else {
		page = *req.Page
	}

	if req.Limit == nil {
		req.Limit = new(int)
		*req.Limit = 10
		limit = *req.Limit
	} else {
		limit = *req.Limit
	}

	var status *string
	if req.Status != nil {
		status = req.Status
	}

	var search *string
	if req.Search != nil {
		search = req.Search
	}

	params := repositories.ListTodoParams{
		Page:     int32(*req.Page),
		LimitVal: int32(*req.Limit),
		Status:   status,
		Search:   search,
	}

	data, err := s.repo.ListTodo(ctx, params)

	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	for _, item := range data {
		todo := Todo{
			ID:          item.ID.String(),
			Title:       item.Title,
			Description: item.Description.String,
			Status:      string(item.Status),
			DueDate:     item.DueDate.Time.Format("2006-01-02"),
		}

		todos = append(todos, todo)
	}

	countData, err = s.repo.CountTodo(ctx, repositories.CountTodoParams{
		Status: status,
		Search: search,
	})

	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	return
}

func (s *TodoService) GetTodo(ctx context.Context, id string) (todo Todo, err error) {

	uuidTodo, err := uuid.Parse(id)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	data, err := s.repo.GetTodoById(ctx, pgtype.UUID{Valid: true, Bytes: uuidTodo})
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	todo = Todo{
		ID:          data.ID.String(),
		Title:       data.Title,
		Description: data.Description.String,
		Status:      string(data.Status),
		DueDate:     data.DueDate.Time.Format("2006-01-02"),
	}

	return

}

func (s *TodoService) UpdateTodo(ctx context.Context, req UpdateTodoRequest, id string) (todo repositories.Todo, err error) {

	timeDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	uuidTodo, err := uuid.Parse(id)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	todo, err = s.repo.UpdateTodo(ctx, repositories.UpdateTodoParams{
		ID:          pgtype.UUID{Valid: true, Bytes: uuidTodo},
		Title:       req.Title,
		Description: pgtype.Text{String: req.Description, Valid: true},
		Status:      repositories.TodoStatus(req.Status),
		DueDate:     pgtype.Date{Time: timeDate, Valid: true},
	})

	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	return
}

func (s *TodoService) DeleteTodo(ctx context.Context, id string) (err error) {

	uuidTodo, err := uuid.Parse(id)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	err = s.repo.DeleteTodo(ctx, pgtype.UUID{Valid: true, Bytes: uuidTodo})

	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	return
}
