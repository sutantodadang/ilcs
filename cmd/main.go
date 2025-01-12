package main

import (
	"context"
	"ilcs/database"
	"ilcs/internal/app/todo"
	"ilcs/internal/http/middlewares"
	"ilcs/internal/http/route"
	"ilcs/internal/repositories"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
}

func main() {

	app := gin.Default()

	db := database.ConnectPG()

	defer db.Close()

	redisDb := database.ConnectRedis()

	defer redisDb.Close()

	app.Use(middlewares.Trace())
	app.Use(middlewares.RequestLoggerMiddleware(), middlewares.ResponseLoggerMiddleware())

	setupContainer(app, db, redisDb)

	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: app,
	}

	go func() {
		log.Info().Msgf("Starting server... on port %s", os.Getenv("PORT"))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")

}

func setupContainer(app *gin.Engine, db *pgxpool.Pool, redisDb *redis.Client) {

	repo := repositories.New(db)

	todoService := todo.NewTodoService(repo, redisDb)

	todoHandler := todo.NewTodoHandler(todoService)

	route.RegisterTodoRoute(app, todoHandler)

}
