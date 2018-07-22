package job

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handler() http.Handler {
	r := chi.NewRouter()

	r.Get("/", h.getAll)

	return r
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
