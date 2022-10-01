package product

//includes a map of ints to products, in a read write mutex, product as key and product as value. multithreaded and maps in go are not thread safe
//thats why we warp our map in mutex to avoid two thread reading and writing at same time
import (
	"database/sql"
	"gowebservice/database"
)

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

func removeProduct(productID int) error {
	_, err := database.DbConn.Query(`DELETE FROM products WHERE productId=?`, productID)
	if err != nil {
		return err
	}
	return nil
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

func updateProduct(product Product) error {
	_, err := database.DbConn.Exec(`UPDATE products SET Manufacturer=?, 
	Sku=?, 
	Upc=?, 
	PricePerUnit=?,
	QuantityOnHand=?,
	ProductName=?`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName,
		product.ProductID) //update the existing product
	if err != nil {
		return err
	}
	return nil
}

func insertProduct(product Product) (int, error) {
	result, err := database.DbConn.Exec(`INSERT INTO products
		(Manufacturer,
		sku,
		upc, 
		PricePerUnit,
		QuantityOnHand,
		ProductName) VALUES (?,?,?,?,?,?)`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName) //inserting into the database using Exec
	if err != nil {
		return 0, nil
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(insertID), nil
}
