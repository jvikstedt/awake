package stepconfig

import (
	"encoding/json"
	"net/http"

	"github.com/jvikstedt/awake/internal/domain"
)

type Handler struct {
	hh                   domain.HandlerHelper
	stepConfigRepository domain.StepConfigRepository
}

func NewHandler(hh domain.HandlerHelper, stepConfigRepository domain.StepConfigRepository) *Handler {
	return &Handler{
		hh:                   hh,
		stepConfigRepository: stepConfigRepository,
	}
}

func (h *Handler) GetOne(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	id, err := h.hh.URLParamInt(r, "id")
	if err != nil {
		return struct{}{}, http.StatusUnprocessableEntity, err
	}

	stepConfig, err := h.stepConfigRepository.GetOne(id)
	if err != nil {
		return struct{}{}, http.StatusInternalServerError, err
	}

	return stepConfig, http.StatusOK, nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	decoder := json.NewDecoder(r.Body)
	stepConfig := domain.StepConfig{}

	if err := decoder.Decode(&stepConfig); err != nil {
		return struct{}{}, http.StatusInternalServerError, err
	}

	newStepConfig, err := h.stepConfigRepository.Create(stepConfig)
	if err != nil {
		return struct{}{}, http.StatusInternalServerError, err
	}

	return newStepConfig, http.StatusOK, nil
}
