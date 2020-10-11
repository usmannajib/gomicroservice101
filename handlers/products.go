package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go-microservice-webinar/data"
	"log"
	"net/http"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

/*func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		rgex := regexp.MustCompile("/([0-9]+)")
		matches := rgex.FindAllStringSubmatch(r.URL.Path, -1)

		if len(matches) != 1 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(matches[0]) != 2 {
			http.Error(rw, "Invalid Ids", http.StatusBadRequest)
			return
		}

		idStr := matches[0][1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(rw, "Id not a num", http.StatusInternalServerError)
			return
		}
		p.updateProduct(id, rw, r)
		return
		//p.l.Println("Got id", id)
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}*/

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	ps := data.GetProducts()
	err := ps.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Error marshalling products", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle HTTP POST")
	product := r.Context().Value(KeyProduct{}).(*data.Product)
	err := product.Validate()
	 if err != nil {
	 	http.Error(rw, fmt.Sprintf("Invalid input JSON: %v", err), http.StatusBadRequest)
	 	return
	 }
	data.AddProduct(product)
	//p.l.Printf("Product: %#v", rp)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to parse id", http.StatusBadRequest)
	}

	p.l.Println("Handle HTTP PUT")
	product := r.Context().Value(KeyProduct{}).(*data.Product)
	data.UpdateProduct(id, product)
}

type KeyProduct struct{}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		product := &data.Product{}
		err := product.FromJSON(request.Body)
		if err != nil {
			http.Error(writer, "Unable to unmarshall product", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(request.Context(), KeyProduct{}, product)
		req := request.WithContext(ctx)
		next.ServeHTTP(writer, req)
	})
}
