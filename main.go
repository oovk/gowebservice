package main

import (
	"gowebservice/database"
	"gowebservice/product"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	product.SetupRoutes(apiBasePath) //
	http.ListenAndServe(":8080", nil)
}
