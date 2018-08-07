package job

import (
	"encoding/json"
	"net/http"

	"github.com/jvikstedt/awake/cron"
	"github.com/jvikstedt/awake/internal/domain"
	"github.com/jvikstedt/awake/internal/runner"
)

type Handler struct {
	hh            domain.HandlerHelper
	jobRepository domain.JobRepository
	runner        *runner.Runner
	scheduler     *cron.Scheduler
}

func NewHandler(hh domain.HandlerHelper, jobRepository domain.JobRepository, runner *runner.Runner, scheduler *cron.Scheduler) *Handler {
	return &Handler{
		hh:            hh,
		jobRepository: jobRepository,
		runner:        runner,
		scheduler:     scheduler,
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

	// if job.Active {
	// 	h.scheduler.AddEntry(cron.EntryID(job.ID), job.Cron, func(id cron.EntryID) {
	// 		h.runner.AddJob(job)
	// 	})
	// } else {
	// 	h.scheduler.RemoveEntry(cron.EntryID(job.ID))
	// }
	h.scheduler.AddEntry(cron.EntryID(job.ID), job.Cron, func(id cron.EntryID) {
		h.runner.AddJob(job)
	})

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

	h.scheduler.AddEntry(cron.EntryID(job.ID), job.Cron, func(id cron.EntryID) {
		h.runner.AddJob(job)
	})

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
