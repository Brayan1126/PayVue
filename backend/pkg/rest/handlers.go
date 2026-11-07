package rest

import (
	"github.com/go-chi/chi/v5"
)

type Handler interface {
	RouteURLs(router *chi.Mux)
}
