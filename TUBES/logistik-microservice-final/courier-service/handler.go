package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"log"
)

type CourierServiceInterface interface {
	StartDelivery(delivery *Delivery) error
	CompleteDelivery(delivery *Delivery) error
	GetCourierDeliveries(deliveries []Delivery, courierID int) []Delivery
	ValidateDelivery(delivery *Delivery) error
}

type CourierHandler struct {
	service CourierServiceInterface
	repository *DeliveryRepository
}

type CompleteDeliveryRequest struct {
    Resi string `json:"resi"`
}

func NewCourierHandler(
	service CourierServiceInterface,
	repository *DeliveryRepository,
) *CourierHandler {

	return &CourierHandler{
		service: service,
		repository: repository,
	}
}

// POST /delivery
func (h *CourierHandler) StartDelivery(w http.ResponseWriter, r *http.Request) {
	var req DeliveryRequest

	// decode request body dari JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	log.Println("COURIER RECEIVED DELIVERY:", req.Resi)

	// validasi field 
	if req.Resi == "" || req.CourierID <= 0 || req.AssignedZone == "" {
		http.Error(w, "resi, courier_id, assigned_zone are required", http.StatusBadRequest)
		return
	}

	delivery := &Delivery{
		Resi:         req.Resi,
		CourierID:    req.CourierID,
		AssignedZone: req.AssignedZone,
		Status:       "pending",
		CreatedAt:    time.Now(),
	}

	// panggil service StartDelivery
	if err := h.service.StartDelivery(delivery); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	err := h.repository.Create(delivery)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// kirim response delivery dalam format JSON
	json.NewEncoder(w).Encode(delivery)
}

// GET /courier/deliveries?courier_id=1
func (h *CourierHandler) GetCourierDeliveries(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("courier_id")

	// validasi query parameter courier_id
	if idStr == "" {
		http.Error(w, "courier_id is required", http.StatusBadRequest)
		return
	}

	courierID, err := strconv.Atoi(idStr)

	// validasi apakah courier_id berupa angka valid
	if err != nil || courierID <= 0 {
		http.Error(w, "invalid courier_id", http.StatusBadRequest)
		return
	}

	// data delivery 
	all := []Delivery{
		{
			Resi:      "", // isi dengan nomor resi
			CourierID: 0,  // isi dengan courier_id
			Status:    "", // isi dengan status delivery
		},
		{
			Resi:      "", // isi dengan nomor resi
			CourierID: 0,  // isi dengan courier_id
			Status:    "", // isi dengan status delivery
		},
	}

	// ambil data delivery berdasarkan courier_id
	result := h.service.GetCourierDeliveries(all, courierID)

	w.Header().Set("Content-Type", "application/json")

	// kirim response hasil delivery courier
	json.NewEncoder(w).Encode(map[string]interface{}{
		"courier_id": courierID,
		"count":      len(result),
		"data":       result,
	})
}

func (h *CourierHandler) CompleteDelivery( w http.ResponseWriter, r *http.Request) {

    var req CompleteDeliveryRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid body", http.StatusBadRequest)
        return
    }

    delivery := &Delivery{
        Resi:   req.Resi,
        Status: "in_delivery",
    }

    if err := h.service.CompleteDelivery(delivery); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(delivery)
}

// GET /health
func (h *CourierHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// response health check service
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy", // isi dengan status service
	})
}
