package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/pelicanch1k/EffectiveMobileTestTask/pkg/logging"
)

const (
	defaultMigrationsDir = "file://migrations"
)

type config struct {
	MigrationsPath string
	DBUsername     string
	DBPassword     string
	DBHost         string
	DBPort         string
	DBName         string
	DBSSLMode      string
}

func main() {
	logger := logging.GetLogger()

	// Определение флагов командной строки
	upCommand := flag.Bool("up", false, "Применить все миграции")
	downCommand := flag.Bool("down", false, "Откатить все миграции")
	stepFlag := flag.Int("step", 0, "Количество шагов для миграции вверх/вниз (0 = все)")
	versionFlag := flag.Bool("version", false, "Показать текущую версию миграции")
	envFileFlag := flag.String("env", ".env", "Путь к .env файлу")
	flag.Parse()

	// Загружаем переменные окружения
	if err := godotenv.Load(*envFileFlag); err != nil && !os.IsNotExist(err) {
		logger.Fatalf("Ошибка загрузки переменных окружения: %s", err)
	}

	// Загружаем конфигурацию
	cfg := loadConfig()

	// Создаем экземпляр миграций
	m, err := migrate.New(cfg.MigrationsPath, buildDatabaseURL(cfg))
	if err != nil {
		logger.Fatalf("Ошибка создания миграции: %v", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			logger.Errorf("Ошибка закрытия источника миграций: %v", srcErr)
		}
		if dbErr != nil {
			logger.Errorf("Ошибка закрытия соединения с БД: %v", dbErr)
		}
	}()

	// Обработка команд
	switch {
	case *versionFlag:
		showMigrationVersion(m, logger)
	case *upCommand:
		applyMigrationsUp(m, *stepFlag, logger)
	case *downCommand:
		applyMigrationsDown(m, *stepFlag, logger)
	default:
		// По умолчанию запускаем миграцию вверх (для обратной совместимости)
		applyMigrationsUp(m, 0, logger)
	}
}

func loadConfig() config {
	return config{
		MigrationsPath: getEnvWithDefault("MIGRATIONS_PATH", defaultMigrationsDir),
		DBUsername:     getEnvWithDefault("DB_USERNAME", "postgres"),
		DBPassword:     getEnvWithDefault("DB_PASSWORD", "postgres"),
		DBHost:         getEnvWithDefault("DB_HOST", "localhost"),
		DBPort:         getEnvWithDefault("DB_PORT", "5432"),
		DBName:         getEnvWithDefault("DB_NAME", "postgres"),
		DBSSLMode:      getEnvWithDefault("DB_SSLMODE", "disable"),
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func buildDatabaseURL(cfg config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)
}

func showMigrationVersion(m *migrate.Migrate, logger *logging.Logger) {
	version, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			logger.Info("База данных не содержит миграций")
			return
		}
		logger.Fatalf("Ошибка получения версии миграции: %v", err)
	}
	logger.Infof("Текущая версия: %d, Состояние: %s", version, dirtyStatus(dirty))
}

func dirtyStatus(dirty bool) string {
	if dirty {
		return "грязная (возможно миграция завершилась с ошибкой)"
	}
	return "чистая"
}

func applyMigrationsUp(m *migrate.Migrate, steps int, logger *logging.Logger) {
	var err error
	if steps > 0 {
		err = m.Steps(steps)
		logger.Infof("Применение %d шагов миграции...", steps)
	} else {
		err = m.Up()
		logger.Info("Применение всех миграций...")
	}

	handleMigrationResult(err, "применения", logger)
}

func applyMigrationsDown(m *migrate.Migrate, steps int, logger *logging.Logger) {
	var err error
	if steps > 0 {
		err = m.Steps(-steps) // Отрицательное значение для шагов вниз
		logger.Infof("Откат %d шагов миграции...", steps)
	} else {
		err = m.Down()
		logger.Info("Откат всех миграций...")
	}

	handleMigrationResult(err, "отката", logger)
}

func handleMigrationResult(err error, operation string, logger *logging.Logger) {
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Infof("Миграции не требуют %s: изменений нет", operation)
		} else {
			logger.Fatalf("Ошибка %s миграций: %v", operation, err)
		}
	} else {
		logger.Infof("Операция %s миграций успешно выполнена!", operation)
	}
}
