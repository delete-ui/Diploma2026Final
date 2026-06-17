package repository

import (
	"GolangBackendDiploma26/internal/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type FavoritesRepository struct {
	db *sql.DB
}

func NewFavoritesRepository(db *sql.DB) *FavoritesRepository {
	return &FavoritesRepository{db: db}
}

func (r *FavoritesRepository) Add(ctx context.Context, userID, batteryID uuid.UUID) error {
	query := `INSERT INTO favorites (user_id, battery_id) 
	          VALUES ($1, $2) 
	          ON CONFLICT (user_id, battery_id) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, userID, batteryID)
	return err
}

func (r *FavoritesRepository) Remove(ctx context.Context, userID, batteryID uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM favorites WHERE user_id = $1 AND battery_id = $2", userID, batteryID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("favorite not found")
	}
	return nil
}

func (r *FavoritesRepository) GetAll(ctx context.Context, userID uuid.UUID) ([]models.Battery, error) {
	query := `SELECT b.id, b.title, b.price, b.stock, b.img, b.brand, 
	          b.voltage, b.polarity, b.capacity, b.standart, b.technology, b.size_type,
	          b.created_at, b.updated_at
	          FROM favorites f 
	          JOIN batteries b ON f.battery_id = b.id 
	          WHERE f.user_id = $1 
	          ORDER BY f.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query favorites: %w", err)
	}
	defer rows.Close()

	var batteries []models.Battery
	for rows.Next() {
		var b models.Battery
		err := rows.Scan(&b.ID, &b.Title, &b.Price, &b.Stock, &b.Img,
			&b.Brand, &b.Voltage, &b.Polarity, &b.Capacity,
			&b.Standart, &b.Technology, &b.SizeType,
			&b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan battery: %w", err)
		}
		batteries = append(batteries, b)
	}
	return batteries, rows.Err()
}

func (r *FavoritesRepository) IsFavorite(ctx context.Context, userID, batteryID uuid.UUID) bool {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id=$1 AND battery_id=$2)", userID, batteryID).Scan(&exists)
	return err == nil && exists
}
