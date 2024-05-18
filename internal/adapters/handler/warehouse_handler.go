package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"warehouse/internal/core/domain"
	"warehouse/internal/core/services"
)

type WarehouseHandler struct {
	svc services.WarehouseService
}

func NewWarehouseHandler(svc services.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{
		svc: svc,
	}
}

func (h *WarehouseHandler) GetWarehouse(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	q := r.URL.Query()

	sID := q.Get("warehouseID")
	if sID == "" {
		ErrorHandler(w, http.StatusBadRequest, errors.New(fmt.Sprintf(errQueryIsEmpty, "warehouseID")))
	}

	ID, err := strconv.Atoi(sID)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, errors.New(fmt.Sprintf(errQueryIsNotNumber, "warehouseID")))
	}

	wh, err := h.svc.GetWarehouse(r.Context(), ID)

	if err != nil {
		if errors.Is(err, services.ErrWarehouseNotFound) {
			ErrorHandler(w, http.StatusBadRequest, err)
			return
		}
		ErrorHandler(w, http.StatusInternalServerError, err)
		return
	}

	SuccessHandler(w, wh)
}

func (h *WarehouseHandler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	wh := domain.Warehouse{}

	if err := json.NewDecoder(r.Body).Decode(&wh); err != nil {
		ErrorHandler(w, http.StatusBadRequest, fmt.Errorf("error decode: %w", err))
		return
	}

	if err := h.svc.CreateWarehouse(r.Context(), wh); err != nil {
		if errors.Is(err, services.ErrGoodIsExist) {
			ErrorHandler(w, http.StatusBadRequest, err)
			return
		}
		ErrorHandler(w, http.StatusInternalServerError, err)
		return
	}

	SuccessHandler(w, wh)
}

func (h *WarehouseHandler) UpdateWarehouse(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	wh := domain.Warehouse{}

	if err := json.NewDecoder(r.Body).Decode(&wh); err != nil {
		ErrorHandler(w, http.StatusBadRequest, fmt.Errorf("error decode: %w", err))
		return
	}

	if err := h.svc.UpdateWarehouse(r.Context(), wh); err != nil {
		if errors.Is(err, services.ErrWarehouseIsExist) {
			ErrorHandler(w, http.StatusBadRequest, err)
			return
		}
		ErrorHandler(w, http.StatusInternalServerError, err)
		return
	}

	SuccessHandler(w, wh)
}

func (h *WarehouseHandler) DeleteWarehouse(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	q := r.URL.Query()

	sID := q.Get("warehouseID")
	if sID == "" {
		ErrorHandler(w, http.StatusBadRequest, errors.New(fmt.Sprintf(errQueryIsEmpty, "warehouseID")))
	}

	ID, err := strconv.Atoi(sID)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, errors.New(fmt.Sprintf(errQueryIsNotNumber, "warehouseID")))
	}

	if err = h.svc.DeleteWarehouse(r.Context(), ID); err != nil {
		if errors.Is(err, services.ErrWarehouseIsNotExist) {
			ErrorHandler(w, http.StatusBadRequest, err)
			return
		}
		ErrorHandler(w, http.StatusInternalServerError, err)
		return
	}

	SuccessHandler(w, nil)
}

func (h *WarehouseHandler) GetCountGoods(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	q := r.URL.Query()

	sID := q.Get("warehouseID")
	if sID == "" {
		ErrorHandler(w, http.StatusBadRequest, errors.New(fmt.Sprintf(errQueryIsEmpty, "warehouseID")))
	}

	ID, err := strconv.Atoi(sID)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, errors.New(fmt.Sprintf(errQueryIsNotNumber, "warehouseID")))
	}

	cnt, err := h.svc.GetCountGoodsByWarehouseID(r.Context(), ID)

	if err != nil {
		if errors.Is(err, services.ErrWarehouseIsNotExist) {
			ErrorHandler(w, http.StatusBadRequest, err)
			return
		}

		ErrorHandler(w, http.StatusInternalServerError, err)
		return
	}

	SuccessHandler(w, map[string]interface{}{
		"count": cnt,
	})
}
