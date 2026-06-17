package repository

import (
	"GolangBackendDiploma26/internal/models"
	"context"
	"database/sql"
	"fmt"
)

type BatteryRepository struct {
	db *sql.DB
}

func NewBatteryRepository(db *sql.DB) *BatteryRepository {
	return &BatteryRepository{db: db}
}

func (r *BatteryRepository) GetAll(ctx context.Context, limit, offset int) ([]models.Battery, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM batteries").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count batteries: %w", err)
	}

	query := `SELECT id, title, price, stock, img, brand, voltage, polarity, capacity, standart, technology, size_type, created_at, updated_at
	          FROM batteries ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query batteries: %w", err)
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
			return nil, 0, fmt.Errorf("scan battery: %w", err)
		}
		batteries = append(batteries, b)
	}
	return batteries, total, rows.Err()
}
