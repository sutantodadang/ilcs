package todo

import (
	"context"
	"encoding/json"
	"errors"
	"ilcs/database"
	"ilcs/internal/constants"
	"ilcs/internal/repositories"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type ITodoService interface {
	CreateTodo(ctx context.Context, req CreateTodoRequest) (todo repositories.Todo, err error)
	GetListTodos(ctx context.Context, req ListTodoRequestParams) (todos []Todo, countData int64, page, limit int, err error)
	GetTodo(ctx context.Context, id string) (todo Todo, err error)
	UpdateTodo(ctx context.Context, req UpdateTodoRequest, id string) (todo repositories.Todo, err error)
	DeleteTodo(ctx context.Context, id string) (err error)
	GetToken(ctx context.Context) (token string, err error)
}

type TodoService struct {
	repo    repositories.Querier
	redisDb database.RedisClient
}

func NewTodoService(repo repositories.Querier, redisDb database.RedisClient) *TodoService {
	return &TodoService{
		repo:    repo,
		redisDb: redisDb,
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

	key := constants.CACHE_KEY + id

	uuidTodo, err := uuid.Parse(id)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	val, err := s.redisDb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {

		data, errG := s.repo.GetTodoById(ctx, pgtype.UUID{Valid: true, Bytes: uuidTodo})
		if errG != nil {
			log.Error().Err(errG).Send()
			err = errG
			return
		}

		todo = Todo{
			ID:          data.ID.String(),
			Title:       data.Title,
			Description: data.Description.String,
			Status:      string(data.Status),
			DueDate:     data.DueDate.Time.Format("2006-01-02"),
		}

		dataByte, errG := json.Marshal(todo)
		if errG != nil {
			log.Error().Err(errG).Send()
			err = errG
			return
		}

		err = s.redisDb.Set(ctx, key, dataByte, 10*time.Minute).Err()
		if err != nil {
			log.Error().Err(err).Send()
			return
		}

	} else if err != nil {
		log.Error().Err(err).Send()
		return
	} else {
		err = json.Unmarshal([]byte(val), &todo)
		if err != nil {
			log.Error().Err(err).Send()
			return

		}
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

func (s *TodoService) GetToken(ctx context.Context) (token string, err error) {

	claimsJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err = claimsJwt.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	return
}
