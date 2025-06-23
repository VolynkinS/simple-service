package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"trainee/internal/api"
	"trainee/internal/config"
	customLogger "trainee/internal/logger"
	"trainee/internal/repo"
	"trainee/internal/service"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func main() {
	if err := godotenv.Load(config.EnvPath); err != nil {
		log.Fatal("Ошибка загрузки env файла:", err)
	}

	// Загружаем конфигурацию из переменных окружения
	var cfg config.AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(errors.Wrap(err, "failed to load configuration"))
	}

	// Инициализация логгера
	logger, err := customLogger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error initializing logger"))
	}

	// Подключение к PostgreSQL
	repository, err := repo.NewRepository(context.Background(), cfg.PostgreSQL)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to initialize repository"))
	}

	// Создание сервиса с бизнес-логикой
	serviceInstance := service.NewService(repository, logger)

	// Инициализация API
	app := api.NewRouters(&api.Routers{Service: serviceInstance}, cfg.Rest.Token)

	// Запуск HTTP-сервера в отдельной горутине
	go func() {
		logger.Infof("Starting server on %s", cfg.Rest.ListenAddress)
		if err := app.Listen(cfg.Rest.ListenAddress); err != nil {
			log.Fatal(errors.Wrap(err, "failed to start server"))
		}
	}()

	// Ожидание системных сигналов для корректного завершения работы
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	logger.Info("Shutting down gracefully...")
}
