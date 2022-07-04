package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	d, err := ioutil.ReadAll(req.Body)
	h.l.Printf("Data %s\n", d)
	// rw.Write(d)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		// rw.WriteHeader(http.StatusBadRequest)
		// rw.Write([]byte("Oops!!!"))
		return
	}

	fmt.Fprintf(rw, "Hellos %s", d)
}
