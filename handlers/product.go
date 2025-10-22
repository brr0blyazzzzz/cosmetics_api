package handlers

import (
	"cosmetics/models"
	"cosmetics/repository"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type ProductHandler struct {
	Repo *repository.ProductRepository
}

// инициализация обработчика
func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{Repo: repo}
}

// извлечение и преобразование данных из HTTP запроса
func parseProductForm(r *http.Request, id int) (*models.Product, error) {
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("ошибка парсинга формы: %w", err)
	}

	title := r.PostFormValue("title")
	description := r.PostFormValue("description")
	application := r.PostFormValue("application")
	photo := r.PostFormValue("photo")

	volumeStr := r.PostFormValue("volume")
	manufacturerIDStr := r.PostFormValue("manufacturer_id")
	volume, err := strconv.ParseFloat(volumeStr, 64)
	if err != nil {
		return nil, fmt.Errorf("Неверный формат объема: %w", err)
	}

	manufacturerID, err := strconv.Atoi(manufacturerIDStr)
	if err != nil {
		return nil, fmt.Errorf("Неверный формат ID производителя: %w", err)
	}
	contraindicationsStr := strings.TrimSpace(r.PostFormValue("contraindications"))
	var contraindications *string
	if contraindicationsStr != "" {
		contraindications = &contraindicationsStr
	}

	return &models.Product{
		ID:                id,
		Title:             title,
		Description:       description,
		Contraindications: contraindications,
		Application:       application,
		Volume:            volume,
		Photo:             photo,
		ManufacturerID:    manufacturerID,
	}, nil
}

// Обработка POST/PUT/DELETE с форм и редирект на админ-панель
func HandleProductFormSubmission(p *ProductHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr, ok := vars["id"]

		var id int
		if ok {
			var err error
			id, err = strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Неверный идентификатор продукта", http.StatusBadRequest)
				return
			}
		}

		r.ParseForm()
		method := r.PostFormValue("_method")

		if method == "DELETE" {
			if err := p.Repo.Delete(id); err != nil {
				log.Printf("Ошибка удаления продукта ID %d: %v", id, err)
				http.Error(w, "Ошибка удаления продукта", http.StatusInternalServerError)
				return
			}
			log.Printf("Успешное удаление продукта ID %d. Редирект на /admin", id)
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}

		if method == "PUT" {
			product, err := parseProductForm(r, id)
			if err != nil {
				log.Printf("Ошибка формы обновления: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if err := p.Repo.Update(product); err != nil {
				log.Printf("Ошибка обновления продукта ID %d: %v", id, err)
				http.Error(w, "Ошибка обновления продукта", http.StatusInternalServerError)
				return
			}
			log.Printf("Успешное обновление продукта ID %d. Редирект на /admin", id)
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}

		if r.Method == "POST" && method == "" {
			product, err := parseProductForm(r, 0)
			if err != nil {
				log.Printf("Ошибка формы создания: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if err := p.Repo.Create(product); err != nil {
				log.Printf("Ошибка создания продукта: %v", err)
				http.Error(w, "Ошибка создания продукта", http.StatusInternalServerError)
				return
			}
			log.Printf("Успешное создание продукта. Редирект на /admin")
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		http.Error(w, "Неверный маршрут", http.StatusMethodNotAllowed)
	}
}

// Обработчик добавления
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.Repo.Create(&product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.Response{Message: "Продукт создан успешно", Data: product})

}

// Обработчик получения всех продуктов
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.Repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{Message: "Продукты получены успешно", Data: products})
}

// Обработчик получения продукта по id
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный идентификатор продукта", http.StatusBadRequest)
		return
	}
	product, err := h.Repo.GetByID(id)
	if err != nil {
		http.Error(w, "Продукт не найден", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{Message: "Продукт получен успешно", Data: product})
}

// Обработчик обновления продукта
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный идентификатор продукта", http.StatusBadRequest)
		return
	}
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	product.ID = id
	if err := h.Repo.Update(&product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{Message: "Продукт обновлен успешно", Data: product})
}

// Обработчик удаления продукта
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный идентификатор продукта", http.StatusBadRequest)
		return
	}
	if err := h.Repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{Message: "Продукт удален успешно"})
}
