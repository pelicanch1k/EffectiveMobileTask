package logging

import (
	"os"
	"path/filepath"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger обертка для zap.Logger с совместимым интерфейсом
type Logger struct {
	zap *zap.SugaredLogger
}

var logger *Logger

// GetLogger возвращает синглтон логгера
func GetLogger() *Logger {
	return logger
}

// Init инициализирует логгер
func Init() {
	// Создаем базовую конфигурацию encoderConfig
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Настраиваем вывод для stdout
	stdout := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	// Создаем директорию для логов, если она не существует
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic(err)
	}

	// Файл для всех логов
	logFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	// Настраиваем вывод для файла
	fileEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileCore := zapcore.NewCore(
		fileEncoder,
		zapcore.AddSync(logFile),
		zapcore.DebugLevel,
	)

	// Объединяем ядра
	core := zapcore.NewTee(stdout, fileCore)

	// Добавляем информацию о вызывающем коде
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// Создаем SugaredLogger для более удобного API
	sugar := zapLogger.Sugar()

	logger = &Logger{
		zap: sugar,
	}

	logger.Info("Zap логгер инициализирован")
}

// Обертки для совместимости с предыдущим интерфейсом

func (l *Logger) Debug(args ...interface{}) {
	l.zap.Debug(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.zap.Debugf(format, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.zap.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.zap.Infof(format, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.zap.Warn(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.zap.Warnf(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.zap.Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.zap.Errorf(format, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.zap.Fatal(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.zap.Fatalf(format, args...)
}

// Функция помощник для получения информации о caller
func getCaller() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown:0"
	}
	return filepath.Base(file) + ":" + string(line)
}

// Явная инициализация логгера (вызывается из main)
func init() {
	// Пустая инициализация, сам логгер будет инициализирован из main
}
