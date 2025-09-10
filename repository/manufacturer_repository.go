package repository

import (
	"cosmetics/models"
	"database/sql"
)

type ManufacturerRepository struct {
	DB *sql.DB
}

// конструктор с подключением
func NewManufacturerRepository(db *sql.DB) *ManufacturerRepository {
	return &ManufacturerRepository{DB: db}
}

// добавление нового производителя
func (r *ManufacturerRepository) Create(manufacturer *models.Manufacturer) error {
	result, err := r.DB.Exec(
		"INSERT INTO manufacturer (manufacturer_title, country, address, contact_list) VALUES (?, ?, ?, ?)",
		manufacturer.Title, manufacturer.Country, manufacturer.Address, manufacturer.ContactList)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	manufacturer.ID = int(id)
	return nil
}

// получение производителя по id
func (r *ManufacturerRepository) GetByID(id int) (*models.Manufacturer, error) {
	var manufacturer models.Manufacturer
	err := r.DB.QueryRow(
		"SELECT manufacturer_id, manufacturer_title, country, address, contact_list FROM manufacturer WHERE manufacturer_id = ?",
		id).Scan(&manufacturer.ID, &manufacturer.Title, &manufacturer.Country, &manufacturer.Address, &manufacturer.ContactList)
	if err != nil {
		return nil, err
	}
	return &manufacturer, nil
}

// получение всех производителей
func (r *ManufacturerRepository) GetAll() ([]models.Manufacturer, error) {
	rows, err := r.DB.Query("SELECT manufacturer_id, manufacturer_title, country, address, contact_list FROM manufacturer")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var manufacturers []models.Manufacturer
	for rows.Next() {
		var manufacturer models.Manufacturer
		if err := rows.Scan(&manufacturer.ID, &manufacturer.Title, &manufacturer.Country, &manufacturer.Address, &manufacturer.ContactList); err != nil {
			return nil, err
		}
		manufacturers = append(manufacturers, manufacturer)
	}
	return manufacturers, nil
}

// обновление произодителя
func (r *ManufacturerRepository) Update(manufacturer *models.Manufacturer) error {
	_, err := r.DB.Exec(
		"UPDATE manufacturer SET manufacturer_title = ?, country = ?, address = ?, contact_list = ? WHERE manufacturer_id = ?",
		manufacturer.Title, manufacturer.Country, manufacturer.Address, manufacturer.ContactList, manufacturer.ID)
	return err
}

// удаление производителя
func (r *ManufacturerRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM manufacturer WHERE manufacturer_id = ?", id)
	return err
}
