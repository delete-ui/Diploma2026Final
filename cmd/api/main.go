package main

import (
	_ "GolangBackendDiploma26/docs"
	"GolangBackendDiploma26/internal/config"
	"GolangBackendDiploma26/internal/logger"
	"GolangBackendDiploma26/internal/server"
	"GolangBackendDiploma26/pkg/postgres"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// @title           Battery Shop API
// @version         1.0
// @description     API для интернет-магазина автомобильных аккумуляторов.
// @description     Предоставляет функциональность для регистрации, аутентификации,
// @description     просмотра каталога аккумуляторов, управления корзиной,
// @description     избранным и оформления заказов(дипломная версия, non-production).
// @description     В связи с дипломной версией сервера все коды подтверждения выводятся в логи или в песочнице аккаунта mailtrap.

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите JWT токен в формате: **Bearer &lt;token&gt;**

// @tag.name auth
// @tag.description Операции аутентификации и управления пользователями

// @tag.name batteries
// @tag.description Операции с каталогом аккумуляторов

// @tag.name cart
// @tag.description Операции с корзиной покупок (требуется авторизация)

// @tag.name favorites
// @tag.description Операции с избранными товарами (требуется авторизация)

// @tag.name orders
// @tag.description Операции с заказами (требуется авторизация)

// @schemes http https
// @accept json
// @produce json
func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Error loading config: %v", err)
		os.Exit(1)
	}

	zapLogger, err := logger.New(logger.Config{
		Level:  cfg.LogLevel,
		Format: cfg.LogFormat,
	})
	if err != nil {
		slog.Error("Error while initializing logger: %v", err)
		os.Exit(1)
	}
	defer zapLogger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pg, err := postgres.New(ctx, &cfg.Database, zapLogger)
	if err != nil {
		slog.Error("Error while connecting to db: %v", err)
		os.Exit(1)
	}
	defer pg.Close()

	zapLogger.Info("application initialized successfully",
		zap.String("env", cfg.Env),
		zap.String("addr", cfg.HTTPServer.Address),
	)

	srv := server.NewServer(*cfg, pg.GetDB(), zapLogger)

	go func() {
		if err := srv.Start(); err != nil {
			zapLogger.Error("server error", zap.Error(err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	zapLogger.Info("application started")

	<-quit
	zapLogger.Info("shutting down gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		zapLogger.Error("server shutdown error", zap.Error(err))
	}

	if err := pg.Close(); err != nil {
		zapLogger.Error("error closing database", zap.Error(err))
	}

	<-shutdownCtx.Done()
	zapLogger.Info("shutdown completed")
}
