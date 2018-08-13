package result

import (
	"encoding/json"
	"net/http"

	"github.com/jvikstedt/awake/internal/domain"
)

type Handler struct {
	resultRepository domain.ResultRepository
}

func NewHandler(resultRepository domain.ResultRepository) *Handler {
	return &Handler{
		resultRepository: resultRepository,
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	results, err := h.resultRepository.GetAll()
	if err != nil {
		return struct{}{}, http.StatusInternalServerError, err
	}

	return results, http.StatusOK, nil
}

func (h *Handler) GetOne(id int, w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	result, err := h.resultRepository.GetOne(id)
	if err != nil {
		return struct{}{}, http.StatusInternalServerError, err
	}

	return result, http.StatusOK, nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	decoder := json.NewDecoder(r.Body)
	result := domain.Result{}

	if err := decoder.Decode(&result); err != nil {
		return struct{}{}, http.StatusInternalServerError, err
	}

	newResult, err := h.resultRepository.Create(result)
	if err != nil {
		return struct{}{}, http.StatusInternalServerError, err
	}

	return newResult, http.StatusOK, nil
}
