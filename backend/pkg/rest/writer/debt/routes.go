package debt

import (
	"github.com/go-chi/chi/v5"
)

func (h *handler) RouteURLs(router *chi.Mux) {
	router.Route("/finances/debt", func(r chi.Router) {
		r.Post("/", h.CreateDebt)
		r.Put("/{id}", h.UpdateDebt)
		r.Delete("/{id}", h.DeleteDebt)
	})
}
