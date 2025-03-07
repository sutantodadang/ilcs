// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CountTodo(ctx context.Context, arg CountTodoParams) (int64, error)
	DeleteTodo(ctx context.Context, id pgtype.UUID) error
	GetTodoById(ctx context.Context, id pgtype.UUID) (GetTodoByIdRow, error)
	InsertTodo(ctx context.Context, arg InsertTodoParams) (Todo, error)
	ListTodo(ctx context.Context, arg ListTodoParams) ([]ListTodoRow, error)
	UpdateTodo(ctx context.Context, arg UpdateTodoParams) (Todo, error)
}

var _ Querier = (*Queries)(nil)
