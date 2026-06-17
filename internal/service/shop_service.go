package service

import (
	"GolangBackendDiploma26/internal/models"
	"GolangBackendDiploma26/internal/repository"
	"GolangBackendDiploma26/internal/transport/email"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ShopService struct {
	cartRepo        *repository.CartRepository
	favoritesRepo   *repository.FavoritesRepository
	transactionRepo *repository.TransactionRepository
	userRepo        *repository.UserRepository
	batteryRepo     *repository.BatteryRepository
	db              *sql.DB
	emailSender     email.Sender
	logger          *zap.Logger
}

func NewShopService(
	cartRepo *repository.CartRepository,
	favoritesRepo *repository.FavoritesRepository,
	transactionRepo *repository.TransactionRepository,
	userRepo *repository.UserRepository,
	batteryRepo *repository.BatteryRepository,
	db *sql.DB,
	emailSender email.Sender,
	logger *zap.Logger,
) *ShopService {
	return &ShopService{
		cartRepo:        cartRepo,
		favoritesRepo:   favoritesRepo,
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		batteryRepo:     batteryRepo,
		db:              db,
		emailSender:     emailSender,
		logger:          logger,
	}
}

func (s *ShopService) AddToCart(ctx context.Context, userID uuid.UUID, batteryID uuid.UUID, quantity int) error {
	return s.cartRepo.AddItem(ctx, userID, batteryID, quantity)
}

func (s *ShopService) GetCart(ctx context.Context, userID uuid.UUID) ([]models.CartItem, error) {
	return s.cartRepo.GetCart(ctx, userID)
}

func (s *ShopService) UpdateCartItem(ctx context.Context, userID, itemID uuid.UUID, quantity int) error {
	return s.cartRepo.UpdateQuantity(ctx, userID, itemID, quantity)
}

func (s *ShopService) RemoveFromCart(ctx context.Context, userID, itemID uuid.UUID) error {
	return s.cartRepo.RemoveItem(ctx, userID, itemID)
}

func (s *ShopService) AddToFavorites(ctx context.Context, userID, batteryID uuid.UUID) error {
	return s.favoritesRepo.Add(ctx, userID, batteryID)
}

func (s *ShopService) RemoveFromFavorites(ctx context.Context, userID, batteryID uuid.UUID) error {
	return s.favoritesRepo.Remove(ctx, userID, batteryID)
}

func (s *ShopService) GetFavorites(ctx context.Context, userID uuid.UUID) ([]models.Battery, error) {
	return s.favoritesRepo.GetAll(ctx, userID)
}

func (s *ShopService) Checkout(ctx context.Context, userID uuid.UUID) (*models.Transaction, error) {
	cartItems, err := s.cartRepo.GetCartForCheckout(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get cart: %w", err)
	}
	if len(cartItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	var total float64
	for _, item := range cartItems {
		total += float64(item.Quantity) * item.Price
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	transaction := &models.Transaction{
		ID:            uuid.New(),
		UserID:        userID,
		TotalAmount:   total,
		Status:        "completed",
		PaymentMethod: "card",
		CreatedAt:     time.Now(),
	}

	if err := s.transactionRepo.Create(ctx, tx, transaction); err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}

	for _, item := range cartItems {
		ti := &models.TransactionItem{
			TransactionID: transaction.ID,
			BatteryID:     item.BatteryID,
			Quantity:      item.Quantity,
			PriceAtTime:   item.Price,
		}
		if err := s.transactionRepo.AddItem(ctx, tx, ti); err != nil {
			return nil, fmt.Errorf("add transaction item: %w", err)
		}
	}

	if _, err := tx.ExecContext(ctx, "DELETE FROM cart WHERE user_id = $1", userID); err != nil {
		return nil, fmt.Errorf("clear cart: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	go s.sendReceipt(user.Email, user.Username, transaction, cartItems)

	return transaction, nil
}

func (s *ShopService) sendReceipt(to, username string, t *models.Transaction, items []models.CartItem) {
	body := fmt.Sprintf(`
		<h2>Заказ #%s</h2>
		<p>Здравствуйте, %s!</p>
		<p>Ваш заказ успешно оформлен.</p>
		<p>Дата: %s</p>
		<p>Сумма: %.2f руб.</p>
		<h3>Товары:</h3>
		<table border="1" cellpadding="5" cellspacing="0">
			<tr><th>Товар</th><th>Количество</th><th>Цена</th><th>Сумма</th></tr>`,
		t.ID.String()[:8], username, t.CreatedAt.Format("02.01.2006 15:04"), t.TotalAmount)

	for _, item := range items {
		body += fmt.Sprintf(`<tr><td>%s</td><td>%d</td><td>%.2f</td><td>%.2f</td></tr>`,
			item.BatteryName, item.Quantity, item.Price, float64(item.Quantity)*item.Price)
	}
	body += `</table><p>Спасибо за покупку!</p>`

	if err := s.emailSender.SendReceipt(to, username, body); err != nil {
		s.logger.Error("failed to send receipt", zap.Error(err))
	}
}

func (s *ShopService) GetOrderHistory(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	transactions, err := s.transactionRepo.GetUserTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i := range transactions {
		items, err := s.transactionRepo.GetTransactionItems(ctx, transactions[i].ID)
		if err != nil {
			return nil, err
		}
		transactions[i].Items = items
	}

	return transactions, nil
}
