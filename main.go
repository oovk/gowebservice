package main

import (
	"gowebservice/database"
	"gowebservice/product"
	"gowebservice/receipt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const basePath = "/api"

func main() {
	database.SetupDatabase()
	receipt.SetupRoutes(basePath)
	product.SetupRoutes(basePath)
	http.ListenAndServe(":8080", nil)
}
