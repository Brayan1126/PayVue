package income

import (
	"github.com/go-chi/chi/v5"
)

func (h *handler) RouteURLs(router *chi.Mux) {
	router.Route("/finances/income", func(r chi.Router) {
		r.Post("/", h.CreateIncome)
		r.Put("/{id}", h.UpdateIncome)
		r.Delete("/{id}", h.DeleteIncome)
	})
}
