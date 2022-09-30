package main

import (
	"gowebservice/product"
	"net/http"
)

const apiBasePath = "/api"

func main() {

	product.SetupRoutes(apiBasePath) //
	http.ListenAndServe(":8080", nil)
}
