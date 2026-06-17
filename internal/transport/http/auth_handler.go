package http

import (
	"GolangBackendDiploma26/internal/models"
	"GolangBackendDiploma26/internal/service"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
)

type AuthHandler struct {
	userService *service.UserService
	validator   *validator.Validate
	logger      *zap.Logger
}

func NewAuthHandler(userService *service.UserService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		validator:   validator.New(),
		logger:      logger,
	}
}

// Register регистрирует нового пользователя
// @Summary      Регистрация нового пользователя
// @Description  Создает учетную запись пользователя и отправляет код подтверждения на указанный email.
// @Description  После регистрации необходимо подтвердить email используя endpoint /verify-email.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.RegisterRequest true "Данные для регистрации"
// @Success      201  {object}  map[string]interface{}  "Пользователь успешно создан"
// @Failure      400  {object}  map[string]string       "Неверный формат запроса"
// @Failure      422  {object}  map[string]string       "Ошибка валидации"
// @Failure      409  {object}  map[string]string       "Пользователь уже существует"
// @Failure      500  {object}  map[string]string       "Внутренняя ошибка сервера"
// @Router       /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode register request", zap.Error(err))
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("validation failed", zap.Error(err))
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	user, err := h.userService.Register(r.Context(), req)
	if err != nil {
		h.logger.Error("registration failed", zap.Error(err))
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "user registered successfully. Check your email for verification code.",
		"user":    user,
	})
}

// VerifyEmail подтверждает email пользователя
// @Summary      Подтверждение email
// @Description  Подтверждает email пользователя с помощью кода, отправленного при регистрации.
// @Description  После успешного подтверждения пользователь может войти в систему.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.VerifyEmailRequest true "Код подтверждения"
// @Success      200  {object}  map[string]string  "Email успешно подтвержден"
// @Failure      400  {object}  map[string]string  "Неверный код или email"
// @Failure      422  {object}  map[string]string  "Ошибка валидации"
// @Failure      500  {object}  map[string]string  "Внутренняя ошибка сервера"
// @Router       /verify-email [post]
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req models.VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode verify request", zap.Error(err))
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("validation failed", zap.Error(err))
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	if err := h.userService.VerifyEmail(r.Context(), req.Email, req.Code); err != nil {
		h.logger.Error("verification failed", zap.Error(err))
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "email verified successfully"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Login авторизует пользователя
// @Summary      Вход в систему
// @Description  Аутентифицирует пользователя по email и паролю.
// @Description  Возвращает JWT токен для доступа к защищенным эндпоинтам.
// @Description  Токен необходимо передавать в заголовке Authorization: Bearer {token}.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "Данные для входа"
// @Success      200  {object}  models.LoginResponse  "Успешный вход"
// @Failure      400  {object}  map[string]string     "Неверный формат запроса"
// @Failure      401  {object}  map[string]string     "Неверный email или пароль"
// @Failure      422  {object}  map[string]string     "Ошибка валидации"
// @Failure      500  {object}  map[string]string     "Внутренняя ошибка сервера"
// @Router       /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode login request", zap.Error(err))
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	resp, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Warn("login failed", zap.Error(err))
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// ForgotPassword инициирует сброс пароля
// @Summary      Забыли пароль
// @Description  Отправляет код для сброса пароля на email пользователя.
// @Description  Код действителен в течение 1 часа.
// @Description  По соображениям безопасности всегда возвращает успех, даже если email не найден.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.ForgotPasswordRequest true "Email пользователя"
// @Success      200  {object}  map[string]string  "Код отправлен (если email существует)"
// @Failure      400  {object}  map[string]string  "Неверный формат запроса"
// @Failure      422  {object}  map[string]string  "Ошибка валидации"
// @Router       /forgot-password [post]
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode forgot password request", zap.Error(err))
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	if err := h.userService.ForgotPassword(r.Context(), req.Email); err != nil {
		h.logger.Error("forgot password error", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "if the email exists, a reset code has been sent"})
}

// ResetPassword сбрасывает пароль пользователя
// @Summary      Сброс пароля
// @Description  Устанавливает новый пароль с использованием кода из письма.
// @Description  После успешного сброса можно войти с новым паролем.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.ResetPasswordRequest true "Данные для сброса пароля"
// @Success      200  {object}  map[string]string  "Пароль успешно изменен"
// @Failure      400  {object}  map[string]string  "Неверный код или email"
// @Failure      422  {object}  map[string]string  "Ошибка валидации"
// @Failure      500  {object}  map[string]string  "Внутренняя ошибка сервера"
// @Router       /reset-password [post]
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode reset password request", zap.Error(err))
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	if err := h.userService.ResetPassword(r.Context(), req.Email, req.Code, req.NewPassword); err != nil {
		h.logger.Warn("reset password failed", zap.Error(err))
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "password has been reset successfully"})
}
