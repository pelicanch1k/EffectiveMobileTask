package app

import (
	"context"
	"net/http"
	"os"
	"time"

	"go.uber.org/fx"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/external_api"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/handler"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/repository"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/repository/postgres"
	router "github.com/pelicanch1k/EffectiveMobileTestTask/internal/router/v1"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/service"
	"github.com/pelicanch1k/EffectiveMobileTestTask/pkg/logging"
)

// loadEnv загружает переменные окружения
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		logging.GetLogger().Fatalf("ошибка загрузки переменных окружения: %s", err.Error())
	}
}

// newLogger инициализирует логгер
func newLogger() *logging.Logger {
	logging.Init()
	return logging.GetLogger()
}

// newDB создает соединение с базой данных
func newDB(logger *logging.Logger) *sqlx.DB {
	db, err := postgres.NewPostgresDB()
	if err != nil {
		logger.Fatalf("ошибка подключения к базе данных: %s", err.Error())
	}
	return db
}

// newServer создает HTTP-сервер
func newServer(r http.Handler, logger *logging.Logger) *http.Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	logger.Infof("Сервер запущен на порту %s", port)

	return srv
}

// registerLifecycle управляет жизненным циклом сервера
func registerLifecycle(lc fx.Lifecycle, srv *http.Server, logger *logging.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatalf("Ошибка запуска сервера: %s", err.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Выключение сервера...")

			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				logger.Fatalf("Ошибка при выключении сервера: %s", err.Error())
			}

			logger.Info("Сервер успешно остановлен")
			return nil
		},
	})
}

func NewApp() *fx.App {
	app := fx.New(
		fx.Invoke(loadEnv),
		fx.Provide(
			newLogger,
			newDB,
			repository.NewRepository,
			fx.Annotate(external_api.NewSongAPI, fx.As(new(service.ExternalAPI))),
			service.NewService,
			handler.NewHandler,
			router.NewRouter,
			fx.Annotate(router.NewRouter, fx.As(new(http.Handler))), 
			newServer,
		),
		fx.Invoke(registerLifecycle),
	)

	return app
}


