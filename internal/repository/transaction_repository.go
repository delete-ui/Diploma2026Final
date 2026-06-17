package repository

import (
	"GolangBackendDiploma26/internal/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *sql.Tx, transaction *models.Transaction) error {
	query := `INSERT INTO transactions (id, user_id, total_amount, status, payment_method) 
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.ExecContext(ctx, query, transaction.ID, transaction.UserID,
		transaction.TotalAmount, transaction.Status, transaction.PaymentMethod)
	return err
}

func (r *TransactionRepository) AddItem(ctx context.Context, tx *sql.Tx, item *models.TransactionItem) error {
	query := `INSERT INTO transaction_items (transaction_id, battery_id, quantity, price_at_time) 
	          VALUES ($1, $2, $3, $4)`
	_, err := tx.ExecContext(ctx, query, item.TransactionID, item.BatteryID,
		item.Quantity, item.PriceAtTime)
	return err
}

func (r *TransactionRepository) GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	query := `SELECT id, user_id, total_amount, status, payment_method, created_at 
	          FROM transactions WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(&t.ID, &t.UserID, &t.TotalAmount, &t.Status,
			&t.PaymentMethod, &t.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}
	return transactions, rows.Err()
}

func (r *TransactionRepository) GetTransactionItems(ctx context.Context, transactionID uuid.UUID) ([]models.TransactionItem, error) {
	query := `SELECT ti.id, ti.transaction_id, ti.battery_id, ti.quantity, ti.price_at_time, ti.created_at,
	          b.title, b.img
	          FROM transaction_items ti 
	          JOIN batteries b ON ti.battery_id = b.id 
	          WHERE ti.transaction_id = $1`
	rows, err := r.db.QueryContext(ctx, query, transactionID)
	if err != nil {
		return nil, fmt.Errorf("query transaction items: %w", err)
	}
	defer rows.Close()

	var items []models.TransactionItem
	for rows.Next() {
		var item models.TransactionItem
		err := rows.Scan(&item.ID, &item.TransactionID, &item.BatteryID,
			&item.Quantity, &item.PriceAtTime, &item.CreatedAt,
			&item.BatteryTitle, &item.BatteryImg)
		if err != nil {
			return nil, fmt.Errorf("scan transaction item: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
