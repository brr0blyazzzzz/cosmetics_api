package main

import (
	"cosmetics/database"
	"cosmetics/handlers"
	"cosmetics/repository"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	//Инициализация базы данных
	if err := database.InitDB(); err != nil {
		log.Fatal("Не удалось подключиться к БД: ", err)
	}

	//Репозитории
	productRepo := repository.NewProductRepository(database.DB)
	manufacturerRepo := repository.NewManufacturerRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)

	//Обработчики
	productHandler := handlers.NewProductHandler(productRepo)
	manufacturerHandler := handlers.NewManufacturerHandler(manufacturerRepo)
	userHandler := handlers.NewUserHandler(userRepo)

	//Маршрутизатор
	r := mux.NewRouter()

	//Статические файлы
	staticFileHandler := http.FileServer(http.Dir("./views"))
	r.PathPrefix("/css/").Handler(staticFileHandler)
	r.PathPrefix("/js/").Handler(staticFileHandler)
	r.PathPrefix("/assets/").Handler(staticFileHandler)

	//Публичные страницы работы с пользователем
	r.HandleFunc("/", handlers.WebHandler(productRepo, manufacturerRepo)).Methods("GET")
	r.HandleFunc("/login", userHandler.LoginPage).Methods("GET")
	r.HandleFunc("/register", userHandler.RegisterPage).Methods("GET")
	r.HandleFunc("/login", userHandler.LoginFormHandler).Methods("POST")
	r.HandleFunc("/register", userHandler.RegisterUser).Methods("POST")
	r.HandleFunc("/logout", handlers.LogoutHandler()).Methods("POST", "GET")

	//Публичные пути продуктов
	r.HandleFunc("/api/products", productHandler.GetProducts).Methods("GET")
	r.HandleFunc("/api/products/{id}", productHandler.GetProduct).Methods("GET")

	//Формы для продукта
	r.HandleFunc("/api/products", handlers.HandleProductFormSubmission(productHandler)).Methods("POST")
	r.HandleFunc("/api/products/{id}", handlers.HandleProductFormSubmission(productHandler)).Methods("POST")

	//Авторизация по JWT-токену
	r.HandleFunc("/api/login", userHandler.LoginUser).Methods("POST")
	r.HandleFunc("/api/register", userHandler.RegisterUser).Methods("POST")

	//Закрытые пути
	api := r.PathPrefix("/api").Subrouter()
	api.Use(handlers.AuthMiddleware)

	api.HandleFunc("/manufacturers", manufacturerHandler.CreateManufacturer).Methods("POST")
	api.HandleFunc("/manufacturers/{id}", manufacturerHandler.UpdateManufacturer).Methods("PUT")
	api.HandleFunc("/manufacturers/{id}", manufacturerHandler.DeleteManufacturer).Methods("DELETE")
	api.HandleFunc("/manufacturers", manufacturerHandler.GetManufacturers).Methods("GET")
	api.HandleFunc("/manufacturers/{id}", manufacturerHandler.GetManufacturer).Methods("GET")

	//Защита админ-панели от неавторизованных пользователей
	r.Handle("/admin", handlers.AuthMiddleware(http.HandlerFunc(handlers.AdminHandler(productRepo)))).Methods("GET")

	//Запуск сервера
	log.Fatal(http.ListenAndServe(":8080", r))
}
