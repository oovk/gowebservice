package main

import "net/http"

type Product struct {
	ProductID    int    `json:"productId"`
	Manufacturer string `json:"manufacturer"`
	Sku          string `json:"sku"`
	Upc          string `json:"upc"`
}

type fooHandler struct {
	Message string
}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(f.Message))
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bar called"))
}

func main() {
	http.Handle("/foo", &fooHandler{Message: "Hellow Vaibhav!"})
	http.HandleFunc("/bar", barHandler)
	http.ListenAndServe(":8080", nil)
}
