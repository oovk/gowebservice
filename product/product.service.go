package product

import (
	"encoding/json"
	"fmt"
	"gowebservice/cors"
	"io/ioutil"
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
	urlPathSegments := strings.Split(r.URL.Path, "products/")               //checking thr productID
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1]) //if it is not integer return error
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	product := getProduct(productID) //if it is a integer the find that in productlist and return the details
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method { // now we check the request typ
	case http.MethodGet: //if it is a get then return the json data from productlist
		productsJSON, err := json.Marshal(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJSON)

	case http.MethodPut: //if the method is put then read the request body and update the product in productlist at perticular index, we are not adding new one just updating details of existing one
		var updateProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updateProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updateProduct.ProductID != productID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		addOrUpdateProduct(updateProduct)
		w.WriteHeader(http.StatusOK)
		return

	case http.MethodDelete: //to handle delete requests
		removeProduct(productID)

	case http.MethodOptions: //cors workflow has preflight request sent by browser using httpOptions method to return the cors specific headers so that browser knows it should allow traffic to be sent to that server
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		producList := getProductList()
		productsJSON, err := json.Marshal(producList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJSON)

	case http.MethodPost:
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newProduct.ProductID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = addOrUpdateProduct(newProduct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated) //201 status code
		return

	case http.MethodOptions: //cors workflow has preflight request sent by browser using httpOptions method to return the cors specific headers so that browser knows it should allow traffic to be sent to that server
		return

	}

}
