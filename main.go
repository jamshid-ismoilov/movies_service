// Package main for movies_service
// @title Movies API
// @version 1.0
// @description A simple movies service with authentication
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"movies_service/auth"
	"movies_service/config"
	"movies_service/handlers"
	"movies_service/repository"
	"movies_service/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	_ "movies_service/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewRouter(userHandler *handlers.UserHandler, movieHandler *handlers.MovieHandler, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	authMiddleware := auth.JWTAuthMiddleware(cfg.JWTSecret)
	movies := router.Group("/movies")
	movies.Use(authMiddleware)
	{
		movies.POST("", movieHandler.CreateMovie)
		movies.GET("", movieHandler.GetMovies)
		movies.GET("/:id", movieHandler.GetMovie)
		movies.PUT("/:id", movieHandler.UpdateMovie)
		movies.DELETE("/:id", movieHandler.DeleteMovie)
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func main() {
	app := fx.New(
		fx.Provide(
			config.NewConfig,
			NewDB,
			repository.NewUserRepository,
			repository.NewMovieRepository,
			func(repo repository.UserRepository, cfg *config.Config) service.UserService {
				return service.NewUserService(repo, cfg.JWTSecret)
			},
			service.NewMovieService,
			handlers.NewUserHandler,
			handlers.NewMovieHandler,
			NewRouter,
			func(lc fx.Lifecycle, router *gin.Engine, cfg *config.Config) *http.Server {
				srv := &http.Server{
					Addr:    ":" + cfg.ServerPort,
					Handler: router,
				}
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
								log.Fatalf("HTTP server error: %v", err)
							}
						}()
						log.Printf("Server running at http://localhost:%s/", cfg.ServerPort)
						return nil
					},
					OnStop: func(ctx context.Context) error {
						log.Println("Shutting down server...")
						return srv.Shutdown(ctx)
					},
				})
				return srv
			},
		),
		fx.Invoke(func(*http.Server) {}),
	)
	app.Run()
}
