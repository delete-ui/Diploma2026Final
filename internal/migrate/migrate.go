package migrate

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"go.uber.org/zap"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Migration struct {
	ID      string
	UpSQL   string
	DownSQL string
}

func LoadMigrations() ([]Migration, error) {
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations dir: %w", err)
	}

	migrationMap := make(map[string]*Migration)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}

		parts := strings.Split(name, ".")
		if len(parts) < 3 {
			continue
		}

		migrationID := parts[0]
		migrationType := parts[1]

		content, err := migrationsFS.ReadFile(fmt.Sprintf("migrations/%s", name))
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", name, err)
		}

		if _, exists := migrationMap[migrationID]; !exists {
			migrationMap[migrationID] = &Migration{ID: migrationID}
		}

		if migrationType == "up" {
			migrationMap[migrationID].UpSQL = string(content)
		} else if migrationType == "down" {
			migrationMap[migrationID].DownSQL = string(content)
		}
	}

	var migrations []Migration
	for _, m := range migrationMap {
		if m.UpSQL == "" {
			return nil, fmt.Errorf("migration %s missing up.sql", m.ID)
		}
		migrations = append(migrations, *m)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].ID < migrations[j].ID
	})

	return migrations, nil
}

func RunMigrations(ctx context.Context, db *sql.DB, logger *zap.Logger) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	migrations, err := LoadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	applied := make(map[string]bool)
	rows, err := db.QueryContext(ctx, "SELECT id FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan migration id: %w", err)
		}
		applied[id] = true
	}

	for _, m := range migrations {
		if applied[m.ID] {
			logger.Debug("migration already applied", zap.String("id", m.ID))
			continue
		}

		logger.Info("applying migration", zap.String("id", m.ID))

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %s: %w", m.ID, err)
		}

		if _, err := tx.ExecContext(ctx, m.UpSQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to apply migration %s: %w\nSQL: %s", m.ID, err, m.UpSQL)
		}

		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (id) VALUES ($1)", m.ID); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", m.ID, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", m.ID, err)
		}

	}

	logger.Info("all migrations completed", zap.Int("count", len(migrations)))
	return nil
}
