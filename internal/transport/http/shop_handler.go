package http

import (
	mw "GolangBackendDiploma26/internal/middleware"
	"GolangBackendDiploma26/internal/models"
	"GolangBackendDiploma26/internal/service"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

type ShopHandler struct {
	shopService *service.ShopService
	logger      *zap.Logger
}

func NewShopHandler(shopService *service.ShopService, logger *zap.Logger) *ShopHandler {
	return &ShopHandler{shopService: shopService, logger: logger}
}

// AddToCart добавляет товар в корзину
// @Summary      Добавить в корзину
// @Description  Добавляет указанный аккумулятор в корзину пользователя.
// @Description  Если товар уже есть в корзине, количество увеличивается.
// @Tags         cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.AddToCartRequest true "ID аккумулятора и количество"
// @Success      200  {object}  map[string]string  "Товар добавлен в корзину"
// @Failure      400  {object}  map[string]string  "Неверный формат запроса"
// @Failure      401  {object}  map[string]string  "Требуется авторизация"
// @Failure      500  {object}  map[string]string  "Внутренняя ошибка сервера"
// @Router       /cart [post]
func (h *ShopHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	user, ok := mw.GetUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req models.AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	userID, _ := uuid.Parse(user.UserID)
	if err := h.shopService.AddToCart(r.Context(), userID, req.BatteryID, req.Quantity); err != nil {
		h.logger.Error("add to cart failed", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "item added to cart"})
}

// GetCart возвращает содержимое корзины
// @Summary      Просмотр корзины
// @Description  Возвращает все товары, находящиеся в корзине текущего пользователя.
// @Tags         cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}  "Содержимое корзины"
// @Failure      401  {object}  map[string]string       "Требуется авторизация"
// @Failure      500  {object}  map[string]string       "Внутренняя ошибка сервера"
// @Router       /cart [get]
func (h *ShopHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	user, ok := mw.GetUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	userID, _ := uuid.Parse(user.UserID)
	items, err := h.shopService.GetCart(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"cart": items})
}

// UpdateCartItem обновляет количество товара в корзине
// @Summary      Обновить количество
// @Description  Изменяет количество указанного товара в корзине.
// @Tags         cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID элемента корзины"
// @Param request body models.UpdateCartItemRequest true "Новое количество"
// @Success      200  {object}  map[string]string  "Корзина обновлена"
// @Failure      400  {object}  map[string]string  "Неверный запрос или товар не найден"
// @Failure      401  {object}  map[string]string  "Требуется авторизация"
// @Failure      500  {object}  map[string]string  "Внутренняя ошибка сервера"
// @Router       /cart/{id} [put]
func (h *ShopHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	user, ok := mw.GetUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	itemID, _ := uuid.Parse(chi.URLParam(r, "id"))
	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	userID, _ := uuid.Parse(user.UserID)
	if err := h.shopService.UpdateCartItem(r.Context(), userID, itemID, req.Quantity); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "cart updated"})
}

// RemoveFromCart удаляет товар из корзины
// @Summary      Удалить из корзины
// @Description  Удаляет указанный товар из корзины пользователя.
// @Tags         cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID элемента корзины"
// @Success      200  {object}  map[string]string  "Товар удален"
// @Failure      400  {object}  map[string]string  "Товар не найден"
// @Failure      401  {object}  map[string]string  "Требуется авторизация"
// @Failure      500  {object}  map[string]string  "Внутренняя ошибка сервера"
// @Router       /cart/{id} [delete]
func (h *ShopHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	user, ok := mw.GetUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	itemID, _ := uuid.Parse(chi.URLParam(r, "id"))
	userID, _ := uuid.Parse(user.UserID)
	if err := h.shopService.RemoveFromCart(r.Context(), userID, itemID); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "item removed"})
}

// Checkout оформляет заказ
// @Summary      Оформить заказ
// @Description  Оформляет заказ на основе содержимого корзины.
// @Description  Корзина очищается, создается запись о транзакции.
// @Description  На email пользователя отправляется квитанция с деталями заказа.
// @Description Оформляет заказ на основе содержимого корзины. Не требует параметров — все товары берутся из корзины текущего пользователя.
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}  "Заказ успешно оформлен"
// @Failure      400  {object}  map[string]string       "Корзина пуста или ошибка оформления"
// @Failure      401  {object}  map[string]string       "Требуется авторизация"
// @Failure      500  {object}  map[string]string       "Внутренняя ошибка сервера"
// @Router       /checkout [post]
func (h *ShopHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	user, ok := mw.GetUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	userID, _ := uuid.Parse(user.UserID)
	transaction, err := h.shopService.Checkout(r.Context(), userID)
	if err != nil {
		h.logger.Error("checkout failed", zap.Error(err))
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message":     "order placed successfully",
		"transaction": transaction,
	})
}

// AddToFavorites добавляет товар в избранное
// @Summary      Добавить в избранное
// @Description  Добавляет указанный аккумулятор в список избранного пользователя.
// @Description  Если товар уже в избранном, запрос игнорируется.
// @Tags         favorites
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param request body models.AddToFavoritesRequest true "ID аккумулятора"
// @Success      200  {object}  map[string]string  "Добавлено в избранное"
// @Failure      400  {object}  map[string]string  "Неверный формат запроса"
// @Failure      401  {object}  map[string]string  "Требуется авторизация"
// @Failure      500  {object}  map[string]string  "Внутренняя ошибка сервера"
// @Router       /favorites [post]
func (h *ShopHandler) AddToFavorites(w http.ResponseWriter, r *http.Request) {
	user, ok := mw.GetUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req struct {
		BatteryID uuid.UUID `json:"battery_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	userID, _ := uuid.Parse(user.UserID)
	if err := h.shopService.AddToFavorites(r.Context(), userID, req.BatteryID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "added to favorites"})
}

// RemoveFromFavorites удаляет товар из избранного
// @Summary      Удалить из избранного
// @Description  Удаляет указанный аккумулятор из списка избранного пользователя.
// @Tags         favorites
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID аккумулятора"
// @Success      200  {object}  map[string]string  "Удалено из избранного"
// @Failure      400  {object}  map[string]string  "Товар не найден в избранном"
// @Failure      401  {object}  map[string]string  "Требуется авторизация"
// @Failure      500  {object}  map[string]string  "Внутренняя ошибка сервера"
// @Router       /favorites/{id} [delete]
func (h *ShopHandler) RemoveFromFavorites(w http.ResponseWriter, r *http.Request) {
	user, ok := mw.GetUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	batteryID, _ := uuid.Parse(chi.URLParam(r, "id"))
	userID, _ := uuid.Parse(user.UserID)
	if err := h.shopService.RemoveFromFavorites(r.Context(), userID, batteryID); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "removed from favorites"})
}

// GetFavorites возвращает список избранного
// @Summary      Просмотр избранного
// @Description  Возвращает все аккумуляторы, добавленные пользователем в избранное.
// @Tags         favorites
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}  "Список избранных аккумуляторов"
// @Failure      401  {object}  map[string]string       "Требуется авторизация"
// @Failure      500  {object}  map[string]string       "Внутренняя ошибка сервера"
// @Router       /favorites [get]
func (h *ShopHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	user, ok := mw.GetUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	userID, _ := uuid.Parse(user.UserID)
	batteries, err := h.shopService.GetFavorites(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"favorites": batteries})
}

// GetOrderHistory возвращает историю заказов
// @Summary      История заказов
// @Description  Возвращает список всех заказов пользователя с деталями каждого заказа.
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}  "Список заказов с товарами"
// @Failure      401  {object}  map[string]string       "Требуется авторизация"
// @Failure      500  {object}  map[string]string       "Внутренняя ошибка сервера"
// @Router       /orders [get]
func (h *ShopHandler) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	user, ok := mw.GetUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	userID, _ := uuid.Parse(user.UserID)
	transactions, err := h.shopService.GetOrderHistory(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"orders": transactions})
}
