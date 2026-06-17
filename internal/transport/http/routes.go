package http

import (
	"GolangBackendDiploma26/internal/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os/exec"
)

func MountRoutes(r chi.Router, auth *AuthHandler, battery *BatteryHandler, shop *ShopHandler, jwtSecret string) {
	authMiddleware := middleware.AuthMiddleware(jwtSecret)

	r.Get("/logs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		cmd := exec.Command("sudo", "docker", "logs", "--tail", "100", "autoshop_api")
		output, err := cmd.Output()
		if err != nil {
			w.Write([]byte("Ошибка получения логов: " + err.Error() + "\n"))
			return
		}

		w.Write(output)
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", auth.Register)
		r.Post("/verify-email", auth.VerifyEmail)
		r.Post("/login", auth.Login)
		r.Post("/forgot-password", auth.ForgotPassword)
		r.Post("/reset-password", auth.ResetPassword)
		r.Get("/batteries", battery.List)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Get("/cart", shop.GetCart)
			r.Post("/cart", shop.AddToCart)
			r.Put("/cart/{id}", shop.UpdateCartItem)
			r.Delete("/cart/{id}", shop.RemoveFromCart)
			r.Post("/checkout", shop.Checkout)
			r.Post("/favorites", shop.AddToFavorites)
			r.Delete("/favorites/{id}", shop.RemoveFromFavorites)
			r.Get("/favorites", shop.GetFavorites)
			r.Get("/orders", shop.GetOrderHistory)
		})
	})
}
