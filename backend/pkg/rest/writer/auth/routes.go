package auth

import (
	"github.com/go-chi/chi/v5"
)

func (h *handler) RouteURLs(router *chi.Mux) {
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/logout", h.Logout)
	})
}
