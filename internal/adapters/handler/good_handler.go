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

var ( //errors
	errQueryIsEmpty     = "query \"%s\" is empty"
	errQueryIsNotNumber = "query \"%s\" is not a number"
)

type GoodHandler struct {
	svc services.GoodService
}

func NewGoodHandler(svc services.GoodService) *GoodHandler {
	return &GoodHandler{
		svc: svc,
	}
}

func (h *GoodHandler) GetGood(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	q := r.URL.Query()
	sID := q.Get("goodID")
	if sID == "" {
		ErrorHandler(w, http.StatusBadRequest, errors.New(fmt.Sprintf(errQueryIsEmpty, "goodID")))
		return
	}
	id, err := strconv.Atoi(sID)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, errors.New(fmt.Sprintf(errQueryIsNotNumber, "goodID")))
	}
	good, err := h.svc.GetGood(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrGoodIDisNegative) || errors.Is(err, services.ErrGoodNotFound) {
			ErrorHandler(w, http.StatusBadRequest, err)
		} else {
			ErrorHandler(w, http.StatusInternalServerError, err)
		}
		return
	}
	SuccessHandler(w, good)
}

func (h *GoodHandler) CreateGood(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	g := domain.Good{}
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		ErrorHandler(w, http.StatusBadRequest, err)
		return
	}
	if err := h.svc.CreateGood(r.Context(), g); err != nil {
		if errors.Is(err, services.ErrGoodIsExist) {
			ErrorHandler(w, http.StatusConflict, err)
		} else {
			ErrorHandler(w, http.StatusInternalServerError, err)
		}
		return
	}
	SuccessHandler(w, g)
}

func (h *GoodHandler) UpdateGood(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	g := domain.Good{}
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		ErrorHandler(w, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.UpdateGood(r.Context(), g); err != nil {
		ErrorHandler(w, http.StatusBadRequest, err)
		return
	}

	SuccessHandler(w, g)
}

func (h *GoodHandler) DeleteGood(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	sID := q.Get("goodID")
	if sID == "" {
		ErrorHandler(w, http.StatusBadRequest, errors.New(fmt.Sprintf(errQueryIsEmpty, "goodID")))
		return
	}

	id, err := strconv.Atoi(sID)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, errors.New(fmt.Sprintf(errQueryIsNotNumber, "goodID")))
		return
	}

	if err = h.svc.DeleteGood(r.Context(), id); err != nil {
		ErrorHandler(w, http.StatusBadRequest, err)
		return
	}

	SuccessHandler(w, nil)
}

func (h *GoodHandler) ReserveGood(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	pairs := make([]domain.PairGoodWarehouse, 0)
	if err := json.NewDecoder(r.Body).Decode(&pairs); err != nil {
		ErrorHandler(w, http.StatusBadRequest, fmt.Errorf("error decode request body: %w", err))
		return
	}

	res, err := h.svc.Reserve(r.Context(), pairs)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, fmt.Errorf("error reserve: %w", err))
		return
	}

	SuccessHandler(w, res)
}

func (h *GoodHandler) ReleaseReservationGood(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	pairs := make([]domain.PairGoodWarehouse, 0)
	if err := json.NewDecoder(r.Body).Decode(&pairs); err != nil {
		ErrorHandler(w, http.StatusBadRequest, fmt.Errorf("error decode request body: %w", err))
		return
	}

	res, err := h.svc.ReleaseReservation(r.Context(), pairs)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, fmt.Errorf("error reserve: %w", err))
		return
	}

	SuccessHandler(w, res)
}

func (h *GoodHandler) AddGoodOnWarehouse(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	q := r.URL.Query()

	sGoodID := q.Get("goodID")
	if sGoodID == "" {
		ErrorHandler(w, http.StatusBadRequest, errors.New(fmt.Sprintf(errQueryIsEmpty, "goodID")))
		return
	}

	goodID, err := strconv.Atoi(sGoodID)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, errors.New(fmt.Sprintf(errQueryIsNotNumber, "goodID")))
		return
	}

	sWarehouseID := q.Get("warehouseID")
	if sWarehouseID == "" {
		ErrorHandler(w, http.StatusBadRequest, errors.New(fmt.Sprintf(errQueryIsEmpty, "warehouseID")))
		return
	}

	warehouseID, err := strconv.Atoi(sWarehouseID)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, errors.New(fmt.Sprintf(errQueryIsNotNumber, "warehouseID")))
		return
	}

	sCount := q.Get("count")

	if sCount == "" {
		ErrorHandler(w, http.StatusBadRequest, errors.New(fmt.Sprintf(errQueryIsEmpty, "count")))
		return
	}

	cnt, err := strconv.Atoi(sCount)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, errors.New(fmt.Sprintf(errQueryIsNotNumber, "count")))
		return
	}

	if err := h.svc.AddGoodOnWarehouse(r.Context(), goodID, warehouseID, cnt); err != nil {
		ErrorHandler(w, http.StatusInternalServerError, err)
		return
	}

	SuccessHandler(w, nil)
}
