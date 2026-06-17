package server

import (
	"GolangBackendDiploma26/internal/models"
	"GolangBackendDiploma26/internal/repository"
	"GolangBackendDiploma26/internal/seed"
	"GolangBackendDiploma26/internal/service"
	"GolangBackendDiploma26/internal/transport/email"
	httpTransport "GolangBackendDiploma26/internal/transport/http"
	"context"
	"database/sql"
	"github.com/go-chi/cors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func NewServer(cfg models.Config, db *sql.DB, logger *zap.Logger) *Server {
	var emailSender email.Sender
	if cfg.SMTP.Host != "" && cfg.SMTP.Username != "" {
		smtpSender, err := email.NewSMTPSender(cfg.SMTP)
		if err != nil {
			logger.Fatal("failed to create SMTP sender", zap.Error(err))
		}
		emailSender = smtpSender
		logger.Info("using SMTP email sender")
	} else {
		emailSender = email.NewLogSender(logger)
		logger.Info("SMTP not configured, using log email sender (stub)")
	}

	userRepo := repository.NewUserRepository(db)
	batteryRepo := repository.NewBatteryRepository(db)
	cartRepo := repository.NewCartRepository(db)
	favoritesRepo := repository.NewFavoritesRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	userService := service.NewUserService(userRepo, logger, cfg.JWT.Secret, cfg.JWT.AccessTTL, emailSender)
	authHandler := httpTransport.NewAuthHandler(userService, logger)
	batteryService := service.NewBatteryService(batteryRepo)
	shopService := service.NewShopService(cartRepo, favoritesRepo, transactionRepo, userRepo, batteryRepo, db, emailSender, logger)
	batteryHandler := httpTransport.NewBatteryHandler(batteryService)
	shopHandler := httpTransport.NewShopHandler(shopService, logger)

	go func() {
		if err := seed.SeedBatteries(context.Background(), db); err != nil {
			logger.Error("failed to seed batteries", zap.Error(err))
		}
	}()

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	httpTransport.MountRoutes(r, authHandler, batteryHandler, shopHandler, cfg.JWT.Secret)

	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.HTTPServer.Address,
			Handler:      r,
			ReadTimeout:  cfg.HTTPServer.ReadTimeout,
			WriteTimeout: cfg.HTTPServer.WriteTimeout,
		},
		logger: logger,
	}
}

func (s *Server) Start() error {
	s.logger.Info("starting HTTP server", zap.String("addr", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down HTTP server")
	return s.httpServer.Shutdown(ctx)
}
