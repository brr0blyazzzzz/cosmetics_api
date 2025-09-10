package models

import "time"

//производитель
type Manufacturer struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Country     string `json:"country"`
	Address     string `json:"address"`
	ContactList string `json:"contact_list"`
}

//состав (единица состава)
type Structure struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//продукт
type Product struct {
	ID                int           `json:"id"`
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	ExpirationDate    time.Time     `json:"expiration_date"`
	Contraindications *string       `json:"contraindications,omitempty"`
	Application       string        `json:"application"`
	Volume            float64       `json:"volume"`
	ManufacturerID    int           `json:"manufacturer_id"`
	Manufacturer      *Manufacturer `json:"manufacturer,omitempty"`
	Structures        []Structure   `json:"structures,omitempty"`
}

//связь многое-ко-многим продукт/единица состава
type ProductStructure struct {
	ProductID   int `json:"product_id"`
	StructureID int `json:"structure_id"`
}

//ответ API
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

//ответ со сведениями об ошибке
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
