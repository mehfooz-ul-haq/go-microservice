package handlers

import (
	"log"
	"net/http"

	"gitlab.com/my-whoosh/admin/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("handle GET Request")

	lp := data.GetProducts()
	// d, err := json.Marshal(lp)
	err := lp.ToJson(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal products json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("handle POST Request")

	prod := &data.Product{}

	err := prod.FromJson(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(prod)
}
