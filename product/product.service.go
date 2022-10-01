package product

import (
	"encoding/json"
	"fmt"
	"gowebservice/cors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const productBasePath = "products"

func SetupRoutes(apiBasePath string) {
	handleProducts := http.HandlerFunc(productsHandler)
	handleProduct := http.HandlerFunc(productHandler)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productBasePath), cors.Middleware(handleProducts)) //makes apipath/products
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productBasePath), cors.Middleware(handleProduct)) //makes apipath/products/
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprint("%s/", productBasePath)) //checking thr productID
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method { // now we check the request typ
	case http.MethodGet: //if it is a get then return the json data from productlist
		product, err := getProduct(productID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if product == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		j, err := json.Marshal(product)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}

	case http.MethodPut: //if the method is put then read the request body and update the product in productlist at perticular index, we are not adding new one just updating details of existing one
		var product Product
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if product.ProductID != productID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = updateProduct(product)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		err := removeProduct(productID)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		producList, err := getProductList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		productsJSON, err := json.Marshal(producList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(productsJSON)
		if err != nil {
			log.Fatal(err)
		}

	case http.MethodPost:
		var newProduct Product
		err := json.NewDecoder(r.Body).Decode(&newProduct)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		productID, err := insertProduct(newProduct)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated) //201 status code
		w.Write([]byte(fmt.Sprintf(`{"productId":%d}`, productID)))
	case http.MethodOptions: //cors workflow has preflight request sent by browser using httpOptions method to return the cors specific headers so that browser knows it should allow traffic to be sent to that server
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
