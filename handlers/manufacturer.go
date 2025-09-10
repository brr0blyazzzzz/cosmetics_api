package handlers

import (
	"cosmetics/models"
	"cosmetics/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ManufacturerHandler struct {
	Repo *repository.ManufacturerRepository
}

// конструктор экземпляра обработчика
func NewManufacturerHandler(repo *repository.ManufacturerRepository) *ManufacturerHandler {
	return &ManufacturerHandler{Repo: repo}
}

// обработчик POST
func (h *ManufacturerHandler) CreateManufacturer(w http.ResponseWriter, r *http.Request) {
	var manufacturer models.Manufacturer
	if err := json.NewDecoder(r.Body).Decode(&manufacturer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.Repo.Create(&manufacturer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.Response{Message: "Производитель создан успешно", Data: manufacturer})
}

// обработчик GETAll
func (h *ManufacturerHandler) GetManufacturers(w http.ResponseWriter, r *http.Request) {
	manufacturers, err := h.Repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{Message: "Производители получены успешно", Data: manufacturers})
}

// обработчик GET
func (h *ManufacturerHandler) GetManufacturer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный идентификатор производителя", http.StatusBadRequest)
		return
	}
	manufacturer, err := h.Repo.GetByID(id)
	if err != nil {
		http.Error(w, "Производитель не найден", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{Message: "Производитель получен успешно", Data: manufacturer})
}

// обработчик PUT
func (h *ManufacturerHandler) UpdateManufacturer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный идентификатор производителя", http.StatusBadRequest)
		return
	}
	var manufacturer models.Manufacturer
	if err := json.NewDecoder(r.Body).Decode(&manufacturer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	manufacturer.ID = id
	if err := h.Repo.Update(&manufacturer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{Message: "Производитель обновлен успешно", Data: manufacturer})
}

// обработчик DELETE
func (h *ManufacturerHandler) DeleteManufacturer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный идентификатор производителя", http.StatusBadRequest)
		return
	}
	if err := h.Repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{Message: "Производитель удален успешно"})
}
