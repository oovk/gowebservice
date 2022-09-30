package product

//includes a map of ints to products, in a read write mutex, product as key and product as value. multithreaded and maps in go are not thread safe
//thats why we warp our map in mutex to avoid two thread reading and writing at same time
import (
	"encoding/json"
	"fmt"
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

func getProduct(ProductID int) *Product {
	productMap.RLock() //read lock to prevent another thread from reading
	defer productMap.RUnlock()
	if product, ok := productMap.m[ProductID]; ok {
		return &product
	}
	return nil
}

func removeProduct(productID int) {
	productMap.Lock() //locking the thread so no other thread can access while we do the action
	defer productMap.Unlock()
	delete(productMap.m, productID)
}

func getProductList() []Product {
	productMap.RLock()                                //read lock untill we read the map
	products := make([]Product, 0, len(productMap.m)) //empty map of product struct type with the same length as productMap
	for _, value := range productMap.m {
		products = append(products, value) //iterating through the productMap and adding the data at index to new product struct type datastructure
	}
	productMap.Unlock() //unlocking after done reading
	return products     //returning product structure
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
		oldProduct := getProduct(product.ProductID)
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
