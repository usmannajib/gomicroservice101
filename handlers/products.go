package handlers

import (
	"go-microservice-webinar/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
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
}

func (p *Products) getProducts(rw http.ResponseWriter) {
	ps := data.GetProducts()
	err := ps.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Error marshalling products", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle HTTP POST")

	rp := &data.Product{}
	err := rp.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshall product", http.StatusBadRequest)
	}

	data.AddProduct(rp)
	//p.l.Printf("Product: %#v", rp)
}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle HTTP PUT")

	rp := &data.Product{}
	err := rp.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshall product", http.StatusBadRequest)
	}

	data.UpdateProduct(id, rp)
}
