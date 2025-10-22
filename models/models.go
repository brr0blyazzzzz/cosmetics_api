package models

import "golang.org/x/crypto/bcrypt"

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

// пользователь
type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"` // В БД будет храниться хеш
}

// Хеширование пароля(используется bcrypt)
func (u *User) SetPassword(password string) error {
	//передача пароля и генерация хеша
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//обработка ошибок
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Сравнение пароля с хешем из БД
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

//продукт
type Product struct {
	ID                int           `json:"id"`
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	Contraindications *string       `json:"contraindications,omitempty"`
	Application       string        `json:"application"`
	Volume            float64       `json:"volume"`
	Photo             string        `json:"photo"`
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
