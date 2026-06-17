package repository

import (
	"GolangBackendDiploma26/internal/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) AddItem(ctx context.Context, userID uuid.UUID, batteryID uuid.UUID, quantity int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO cart (user_id, battery_id, quantity) 
	          VALUES ($1, $2, $3) 
	          ON CONFLICT (user_id, battery_id) 
	          DO UPDATE SET quantity = cart.quantity + $3, updated_at = NOW()`
	_, err = tx.ExecContext(ctx, query, userID, batteryID, quantity)
	if err != nil {
		return fmt.Errorf("insert/update cart: %w", err)
	}

	return tx.Commit()
}

func (r *CartRepository) GetCart(ctx context.Context, userID uuid.UUID) ([]models.CartItem, error) {
	query := `SELECT c.id, c.user_id, c.battery_id, c.quantity, c.created_at, c.updated_at,
	          b.title, b.price, b.img, b.brand
	          FROM cart c 
	          JOIN batteries b ON c.battery_id = b.id 
	          WHERE c.user_id = $1 
	          ORDER BY c.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query cart: %w", err)
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		err := rows.Scan(&item.ID, &item.UserID, &item.BatteryID, &item.Quantity,
			&item.CreatedAt, &item.UpdatedAt, &item.BatteryName, &item.Price,
			&item.Img, &item.Brand)
		if err != nil {
			return nil, fmt.Errorf("scan cart item: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *CartRepository) UpdateQuantity(ctx context.Context, userID, itemID uuid.UUID, quantity int) error {
	query := `UPDATE cart SET quantity = $1, updated_at = NOW() 
	          WHERE id = $2 AND user_id = $3`
	result, err := r.db.ExecContext(ctx, query, quantity, itemID, userID)
	if err != nil {
		return fmt.Errorf("update quantity: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("item not found")
	}
	return nil
}

func (r *CartRepository) RemoveItem(ctx context.Context, userID, itemID uuid.UUID) error {
	query := `DELETE FROM cart WHERE id = $1 AND user_id = $2`
	result, err := r.db.ExecContext(ctx, query, itemID, userID)
	if err != nil {
		return fmt.Errorf("delete cart item: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("item not found")
	}
	return nil
}

func (r *CartRepository) ClearCart(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM cart WHERE user_id = $1", userID)
	return err
}

func (r *CartRepository) GetCartForCheckout(ctx context.Context, userID uuid.UUID) ([]models.CartItem, error) {
	return r.GetCart(ctx, userID)
}
