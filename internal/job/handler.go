package job

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jvikstedt/awake/internal/domain"
)

type Handler struct {
	jobRepository domain.JobRepository
}

func NewHandler(jobRepository domain.JobRepository) *Handler {
	return &Handler{
		jobRepository: jobRepository,
	}
}

func (h *Handler) Handler() http.Handler {
	r := chi.NewRouter()

	r.Get("/", h.getAll)
	r.Get("/{id}", h.getOne)

	return r
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.jobRepository.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(jobs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getOne(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	job, err := h.jobRepository.GetOne(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(job)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
