package main

import (
	"cosmetics/database"
	"cosmetics/handlers"
	"cosmetics/repository"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	//инициализация и подключение к базе данных
	if err := database.InitDB(); err != nil {
		log.Fatal("Не удалось подключиться к базе косметических продуктов :(", err)
	}

	//инициализация репозиториев
	productRepo := repository.NewProductRepository(database.DB)
	manufacturerRepo := repository.NewManufacturerRepository(database.DB)

	//инициализация обработчиков
	productHandler := handlers.NewProductHandler(productRepo)
	manufacturerHandler := handlers.NewManufacturerHandler(manufacturerRepo)

	r := mux.NewRouter()

	//маршруты для продуктов
	r.HandleFunc("/api/products", productHandler.GetProducts).Methods("GET")
	r.HandleFunc("/api/products/{id}", productHandler.GetProduct).Methods("GET")
	r.HandleFunc("/api/products", productHandler.CreateProduct).Methods("POST")
	r.HandleFunc("/api/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	r.HandleFunc("/api/products/{id}", productHandler.DeleteProduct).Methods("DELETE")

	//маршруты для производителей
	r.HandleFunc("/api/manufacturers", manufacturerHandler.GetManufacturers).Methods("GET")
	r.HandleFunc("/api/manufacturers/{id}", manufacturerHandler.GetManufacturer).Methods("GET")
	r.HandleFunc("/api/manufacturers", manufacturerHandler.CreateManufacturer).Methods("POST")
	r.HandleFunc("/api/manufacturers/{id}", manufacturerHandler.UpdateManufacturer).Methods("PUT")
	r.HandleFunc("/api/manufacturers/{id}", manufacturerHandler.DeleteManufacturer).Methods("DELETE")

	//главная страница
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Добро пожаловать в API косметических продуктов :)"})
	}).Methods("GET")

	log.Println("Сервер :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
