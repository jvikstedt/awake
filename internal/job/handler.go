package job

import (
	"encoding/json"
	"net/http"

	"github.com/jvikstedt/awake/internal/domain"
)

type Handler struct {
	hh            domain.HandlerHelper
	jobRepository domain.JobRepository
}

func NewHandler(hh domain.HandlerHelper, jobRepository domain.JobRepository) *Handler {
	return &Handler{
		hh:            hh,
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
	id, err := h.hh.URLParamInt(r, "id")
	if err != nil {
		return struct{}{}, http.StatusUnprocessableEntity, err
	}

	job, err := h.jobRepository.GetOne(id)
	if err != nil {
		return struct{}{}, http.StatusInternalServerError, err
	}

	return job, http.StatusOK, nil
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	id, err := h.hh.URLParamInt(r, "id")
	if err != nil {
		return struct{}{}, http.StatusUnprocessableEntity, err
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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	decoder := json.NewDecoder(r.Body)
	job := domain.Job{}

	if err := decoder.Decode(&job); err != nil {
		return struct{}{}, http.StatusInternalServerError, err
	}

	newJob, err := h.jobRepository.Create(job)
	if err != nil {
		return struct{}{}, http.StatusInternalServerError, err
	}

	return newJob, http.StatusOK, nil
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	id, err := h.hh.URLParamInt(r, "id")
	if err != nil {
		return struct{}{}, http.StatusUnprocessableEntity, err
	}

	newJob, err := h.jobRepository.Delete(id)
	if err != nil {
		return struct{}{}, http.StatusInternalServerError, err
	}

	return newJob, http.StatusOK, nil
}
