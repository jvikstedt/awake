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

func (h *Handler) Routes(r chi.Router) {
	r.Get("/", h.getAll)
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
