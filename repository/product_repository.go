package repository

import (
	"cosmetics/models"
	"database/sql"
)

type ProductRepository struct {
	DB *sql.DB
}

// конструктор с подключением
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

// добавление нового продукта
func (r *ProductRepository) Create(product *models.Product) error {
	result, err := r.DB.Exec(`INSERT INTO products (product_title, product_description, expiration_date, contraindications, application, volume, manufacturer_id) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		product.Title, product.Description, product.ExpirationDate, product.Contraindications, product.Application, product.Volume, product.ManufacturerID)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	product.ID = int(id)
	return nil
}

// получение продукта по id
func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	var product models.Product
	var contraindications sql.NullString

	err := r.DB.QueryRow(`SELECT product_id, product_title, product_description, expiration_date, contraindications, application, volume, manufacturer_id FROM products WHERE product_id = ?`,
		id).Scan(&product.ID, &product.Title, &product.Description, &product.ExpirationDate, &contraindications, &product.Application, &product.Volume, &product.ManufacturerID)
	if err != nil {
		return nil, err
	}
	if contraindications.Valid {
		product.Contraindications = &contraindications.String
	}

	manufacturer, _ := NewManufacturerRepository(r.DB).GetByID(product.ManufacturerID)
	product.Manufacturer = manufacturer

	structureRows, err := r.DB.Query(`SELECT s.structure_id, s.structure_name FROM structure s JOIN product_structure ps ON ps.structure_id = s.structure_id WHERE ps.product_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer structureRows.Close()

	var structures []models.Structure
	for structureRows.Next() {
		var s models.Structure
		if err := structureRows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}
		structures = append(structures, s)
	}
	product.Structures = structures
	return &product, nil
}

// получение всех продуктов
func (r *ProductRepository) GetAll() ([]models.Product, error) {
	rows, err := r.DB.Query(`SELECT product_id, product_title, product_description, expiration_date, contraindications, application, volume, manufacturer_id FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		var contraindications sql.NullString

		if err := rows.Scan(&product.ID, &product.Title, &product.Description, &product.ExpirationDate, &contraindications, &product.Application, &product.Volume, &product.ManufacturerID); err != nil {
			return nil, err
		}
		if contraindications.Valid {
			product.Contraindications = &contraindications.String
		}

		manufacturer, _ := NewManufacturerRepository(r.DB).GetByID(product.ManufacturerID)
		product.Manufacturer = manufacturer

		structureRows, err := r.DB.Query(`SELECT s.structure_id, s.structure_name FROM structure s JOIN product_structure ps ON ps.structure_id = s.structure_id WHERE ps.product_id = ?`, product.ID)
		if err != nil {
			return nil, err
		}
		defer structureRows.Close()

		var structures []models.Structure
		for structureRows.Next() {
			var s models.Structure
			if err := structureRows.Scan(&s.ID, &s.Name); err != nil {
				return nil, err
			}
			structures = append(structures, s)
		}
		product.Structures = structures

		products = append(products, product)
	}
	return products, nil
}

// обновление продукта
func (r *ProductRepository) Update(product *models.Product) error {
	_, err := r.DB.Exec(`UPDATE products SET product_title = ?, product_description = ?, expiration_date = ?, contraindications = ?, application = ?, volume = ?, manufacturer_id = ? WHERE product_id = ?`,
		product.Title, product.Description, product.ExpirationDate, product.Contraindications, product.Application, product.Volume, product.ManufacturerID, product.ID)
	return err
}

// удаление продукта
func (r *ProductRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM products WHERE product_id = ?", id)
	return err
}
