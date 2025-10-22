package repository

import (
	"cosmetics/models"
	"database/sql"
	"fmt"
	"log"
	"strings"
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
	result, err := r.DB.Exec(`INSERT INTO products (product_title, product_description,contraindications, application, volume, manufacturer_id, photo) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		product.Title, product.Description, product.Contraindications, product.Application, product.Volume, product.ManufacturerID, product.Photo)
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
	err := r.DB.QueryRow(`SELECT product_id, product_title, product_description, contraindications, application, volume, photo, manufacturer_id FROM products WHERE product_id = ?`,
		id).Scan(&product.ID, &product.Title, &product.Description, &contraindications, &product.Application, &product.Volume, &product.Photo, &product.ManufacturerID)
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
	rows, err := r.DB.Query(`SELECT product_id, product_title, product_description, contraindications, application, volume, photo, manufacturer_id FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		var contraindications sql.NullString

		if err := rows.Scan(&product.ID, &product.Title, &product.Description, &contraindications, &product.Application, &product.Volume, &product.Photo, &product.ManufacturerID); err != nil {
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
	_, err := r.DB.Exec(`UPDATE products SET product_title = ?, product_description = ?, contraindications = ?, application = ?, volume = ?, photo = ?, manufacturer_id = ? WHERE product_id = ?`,
		product.Title, product.Description, product.Contraindications, product.Application, product.Volume, product.Photo, product.ManufacturerID, product.ID)
	return err
}

// удаление продукта
func (r *ProductRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM products WHERE product_id = ?", id)
	return err
}

// получение продукта по заданным требованиями(по названию, по производителю)
func (r *ProductRepository) GetProductsSearch(manufacturerID int, searchQuery string) ([]models.Product, error) {
	var products []models.Product
	var args []interface{}
	argCount := 0
	query := `
        SELECT p.product_id, p.product_title, p.product_description, p.contraindications, p.application, p.volume, p.photo, p.manufacturer_id,
               m.manufacturer_id, m.manufacturer_title
        FROM products p
        JOIN manufacturer m ON p.manufacturer_id = m.manufacturer_id
    `
	whereClauses := []string{}
	if manufacturerID > 0 {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("p.manufacturer_id = $%d", argCount))
		args = append(args, manufacturerID)
	}
	if searchQuery != "" {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("p.product_title LIKE $%d", argCount))
		args = append(args, "%"+searchQuery+"%")
	}
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	query += " ORDER BY p.product_id ASC"
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		log.Printf("Ошибка выполнения запроса с фильтрами: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Product
		var m models.Manufacturer
		var contraindications sql.NullString
		err := rows.Scan(
			&p.ID, &p.Title, &p.Description, &contraindications, &p.Application, &p.Volume, &p.Photo, &p.ManufacturerID,
			&m.ID, &m.Title,
		)
		if err != nil {
			log.Printf("Ошибка продукта: %v", err)
			return nil, err
		}

		if contraindications.Valid {
			p.Contraindications = &contraindications.String
		}
		p.Manufacturer = &m
		products = append(products, p)
	}

	return products, nil
}
