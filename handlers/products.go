package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"gitlab.com/my-whoosh/admin/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	// get all products
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	// save new product
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	//update a product

	if r.Method == http.MethodPut {
		// expect the id in the URI
		p.l.Println("PUT request", r.URL.Path)

		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
		}

		if len(g[0]) != 1 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)

		if err != nil {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
		}

		p.l.Println("got id", id)

		p.updateProduct(id, rw, r)

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

	// p.l.Printf("Prod: %#v", prod)
	data.AddProduct(prod)
}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	prod := &data.Product{}

	err := prod.FromJson(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err != data.UpdateProduct(prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return 
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return 
	}
}
