package postgres

import (
	"GolangBackendDiploma26/internal/config"
	"GolangBackendDiploma26/internal/migrate"
	"GolangBackendDiploma26/internal/models"
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"time"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db     *sql.DB
	logger *zap.Logger
	config *models.DatabaseConfig
}

type Config struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func New(ctx context.Context, cfg *models.DatabaseConfig, logger *zap.Logger) (*Postgres, error) {
	dsn := config.DSN(*cfg)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("database connected",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
		zap.String("dbname", cfg.DBName),
	)

	pg := &Postgres{
		db:     db,
		logger: logger,
		config: cfg,
	}

	if err := pg.applyMigrations(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return pg, nil
}

func (p *Postgres) applyMigrations(ctx context.Context) error {
	p.logger.Info("running database migrations...")
	if err := migrate.RunMigrations(ctx, p.db, p.logger); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	p.logger.Info("database migrations completed")
	return nil
}

func (p *Postgres) GetDB() *sql.DB {
	return p.db
}

func (p *Postgres) Close() error {
	p.logger.Info("closing database connection")
	return p.db.Close()
}

func (p *Postgres) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

func (p *Postgres) InTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			p.logger.Error("failed to rollback transaction", zap.String("error", rbErr.Error()))
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
