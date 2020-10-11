package data

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"io"
	"regexp"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU, true)
	/*err :=
	if err != nil {
		return fmt.Errorf("Invalid SKU")
	}*/
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	ex := regexp.MustCompile("[a-z]+-[a-z]+-[a-z]+")
	matches := ex.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
}

func GetProducts() Products {
	return productsList
}

func AddProduct(p *Product) {
	(*p).ID = getNextID()
	p.ID = getNextID()
	productsList = append(productsList, p)
}

func UpdateProduct(id int, p *Product) error {
	_, i, err := findProduct(id)
	if err != nil {
		return err
	}

	p.ID = id
	productsList[i] = p
	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found.")

func findProduct(id int) (*Product, int, error) {
	for i, p := range productsList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, 0, ErrProductNotFound
}

func getNextID() int {
	return productsList[len(productsList)-1].ID + 1
}

var productsList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
