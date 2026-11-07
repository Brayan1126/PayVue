package payment

import (
	"github.com/go-chi/chi/v5"
)

func (h *handler) RouteURLs(router *chi.Mux) {
	router.Route("/finances/payment", func(r chi.Router) {
		r.Get("/", h.GetAllPayments)
		r.Get("/receipt/{filename}", h.GetReceipt)
	})
}
