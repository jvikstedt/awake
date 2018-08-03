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

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	jobs, err := h.jobRepository.GetAll()
	if err != nil {
		return struct{}{}, http.StatusInternalServerError, err
	}

	return jobs, http.StatusOK, nil
}

func (h *Handler) GetOne(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return struct{}{}, http.StatusInternalServerError, err
	}

	job, err := h.jobRepository.GetOne(id)
	if err != nil {
		return struct{}{}, http.StatusInternalServerError, err
	}

	return job, http.StatusOK, nil
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return struct{}{}, http.StatusInternalServerError, err
	}

	decoder := json.NewDecoder(r.Body)
	job := domain.Job{}

	if err := decoder.Decode(&job); err != nil {
		return struct{}{}, http.StatusInternalServerError, err
	}

	newJob, err := h.jobRepository.Update(id, job)
	if err != nil {
		return struct{}{}, http.StatusInternalServerError, err
	}

	return newJob, http.StatusOK, nil
}
