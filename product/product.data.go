package product

//includes a map of ints to products, in a read write mutex, product as key and product as value. multithreaded and maps in go are not thread safe
//thats why we warp our map in mutex to avoid two thread reading and writing at same time
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gowebservice/database"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
)

var productMap = struct {
	sync.RWMutex
	m map[int]Product
}{m: make(map[int]Product)}

func init() {
	fmt.Println("Loading Products....")
	prodMap, err := loadProductMap()
	productMap.m = prodMap
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d products loaded...\n", len(productMap.m))
}

func loadProductMap() (map[int]Product, error) {
	fileName := "products.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	} //file does not exist return an error checked by using os.stat package

	file, _ := ioutil.ReadFile(fileName) //read the byte data
	productList := make([]Product, 0)
	err = json.Unmarshal([]byte(file), &productList) //deserialize the file data into slice of products
	if err != nil {
		log.Fatal(err)
	}
	prodMap := make(map[int]Product) //iterate over the slice of products to initialize the item in our map
	for i := 0; i < len(productList); i++ {
		prodMap[productList[i].ProductID] = productList[i]
	}
	return prodMap, nil

}

func getProduct(ProductID int) (*Product, error) {
	row := database.DbConn.QueryRow(`SELECT productId,
	Manufaturer, 
	sku, 
	upc,
	pricePerUnit,
	quantityOnHand,
	productName, 
	FROM products
	WHERE productId = ?`, ProductID) //get the specifc row corresponding to productID
	product := &Product{}
	err := row.Scan(&product.ProductID,
		&product.Manufacturer,
		&product.Sku,
		&product.Upc,
		&product.PricePerUnit,
		&product.QuantityOnHand,
		&product.ProductName)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return product, nil
}

func removeProduct(productID int) {
	productMap.Lock() //locking the thread so no other thread can access while we do the action
	defer productMap.Unlock()
	delete(productMap.m, productID)
}

func getProductList() ([]Product, error) {
	results, err := database.DbConn.Query(`SELECT productId,
	Manufaturer, 
	sku, 
	upc,
	pricePerUnit,
	quantityOnHand,
	productName, 
	FROM products`)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	products := make([]Product, 0) // make slice for product mapping with productId
	for results.Next() {           //move the product to scan next record
		var product Product
		results.Scan(&product.ProductID,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.PricePerUnit,
			&product.QuantityOnHand,
			&product.ProductName)
		products = append(products, product)
	}
	return products, nil
}

func getProductIds() []int {
	productMap.RLock()
	productIds := []int{} //list of integer values
	for key := range productMap.m {
		productIds = append(productIds, key) //appending all the key values to productIds
	}
	productMap.RUnlock()
	sort.Ints(productIds) //sort the slice of int in ascending order
	return productIds
}

func getNextProductID() int { //get the next product id where we need to append for example if list has 4 items this function will return 5 and we'll add to 5th position
	productIDs := getProductIds()
	return productIDs[len(productIDs)-1] + 1
}

func addOrUpdateProduct(product Product) (int, error) {
	addOrUpdateID := -1
	if product.ProductID > 0 {
		oldProduct, err := getProduct(product.ProductID)
		if err != nil {
			return addOrUpdateID, err
		}
		if oldProduct == nil {
			return 0, fmt.Errorf("product id [%d] does not exists", product.ProductID)
		}
		addOrUpdateID = product.ProductID //modify the existing items in the map
	} else { //used when item is not present in the map and we are adding new item to the map
		addOrUpdateID = getNextProductID() //getting the index where we need to add
		product.ProductID = addOrUpdateID  //setting the product id for product information we need to add
	}
	productMap.Lock()
	productMap.m[addOrUpdateID] = product //updating the productMap with new details or adding new details to map
	productMap.Unlock()
	return addOrUpdateID, nil

}
