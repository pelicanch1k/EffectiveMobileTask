package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/handler"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/repository"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/repository/postgres"
	router "github.com/pelicanch1k/EffectiveMobileTestTask/internal/router/v1"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/service"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/external_api"
	"github.com/pelicanch1k/EffectiveMobileTestTask/pkg/logging"
)

// @title Song API
// @version 1.0
// @description API Server for EffectiveMobileTestTask

// @host localhost:80
func main() {
	// TODO: init logger
	logging.Init()
	logger := logging.GetLogger()

	// TODO: init .env
	if err := godotenv.Load(); err != nil {
		logger.Fatalf("ошибка загрузки переменных окружения: %s", err.Error())
	}

	// TODO: init db
	db, err := postgres.NewPostgresDB()
	if err != nil {
		logger.Fatalf("ошибка подключения к базе данных: %s", err.Error())
	}

	// TODO: init main components
	repos := repository.NewRepository(db)
	externalAPI := external_api.NewSongAPI()
	services := service.NewService(repos, externalAPI)
	handlers := handler.NewHandler(services, logger)

	// TODO: init router
	r := router.NewRouter(handlers)

	// Настройка сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // порт по умолчанию
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Run server
	go func() {
		logger.Infof("Запуск сервера на порту %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("ошибка запуска сервера: %s", err.Error())
		}
	}()

	// Graceful
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Выключение сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("ошибка при выключении сервера: %s", err.Error())
	}

	logger.Info("Сервер успешно остановлен")
}
