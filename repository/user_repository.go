package repository

import (
	"cosmetics/models"
	"database/sql"
	"fmt"
	"log"
)

type UserRepository struct {
	DB *sql.DB
}

// конструктор с подключением
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// создание пользователя
func (r *UserRepository) CreateUser(user *models.User) error {
	if r.DB == nil {
		return fmt.Errorf("database connection is nil")
	}
	query := "INSERT INTO users (username, password) VALUES (?, ?)"

	result, err := r.DB.Exec(query, user.UserName, user.Password)
	if err != nil {
		log.Printf("DB Error (CreateUser): %v", err)
		return fmt.Errorf("ошибка при создании пользователя: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("пользователь не был создан")
	}

	fmt.Printf("--- Пользователь '%s' успешно сохранен в БД ---\n", user.UserName)
	return nil
}

// получение пользователя по логину
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	user := &models.User{}
	query := "SELECT id, username, password FROM users WHERE username = ?"

	row := r.DB.QueryRow(query, username)

	err := row.Scan(&user.ID, &user.UserName, &user.Password)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("DB Error (GetUserByUsername): %v", err)
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}

	return user, nil
}
