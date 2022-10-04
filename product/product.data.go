package product

//includes a map of ints to products, in a read write mutex, product as key and product as value. multithreaded and maps in go are not thread safe
//thats why we warp our map in mutex to avoid two thread reading and writing at same time
import (
	"context"
	"database/sql"
	"errors"
	"gowebservice/database"
	"log"
	"strings"
	"time"
)

func getProduct(ProductID int) (*Product, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 15*time.Second) //if the query is going to take longer than 15 sec then its going to cancle and return
	defer cancle()
	row := database.DbConn.QueryRowContext(ctx, `SELECT productId,
	manufacturer, 
	sku, 
	upc,
	pricePerUnit,
	quantityOnHand,
	productName 
	FROM products
	WHERE productId=?`, ProductID) //get the specifc row corresponding to productID
	product := &Product{}
	err := row.Scan(
		&product.ProductID,
		&product.Manufacturer,
		&product.Sku,
		&product.Upc,
		&product.PricePerUnit,
		&product.QuantityOnHand,
		&product.ProductName,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return product, nil
}

func removeProduct(productID int) error {
	ctx, cancle := context.WithTimeout(context.Background(), 15*time.Second) //if the query is going to take longer than 15 sec then its going to cancle and return
	defer cancle()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM products WHERE productId = ?`, productID)
	if err != nil {
		log.Print(err.Error())
		return err
	}
	return nil
}

func getProductList() ([]Product, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 15*time.Second) //if the query is going to take longer than 15 sec then its going to cancle and return
	defer cancle()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
	productId,
	manufacturer, 
	sku, 
	upc,
	pricePerUnit,
	quantityOnHand,
	productName 
	FROM products`)
	if err != nil {
		log.Println(err.Error())
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

func GetTopTenProducts() ([]Product, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 15*time.Second) //if the query is going to take longer than 15 sec then its going to cancle and return
	defer cancle()
	results, err := database.DbConn.QueryContext(ctx, `SELECT
	productId,
	manufacturer, 
	sku, 
	upc, 
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products ORDER BY quantityOnHand DESC LIMIT 10
	`)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	defer results.Close()
	products := make([]Product, 0)
	for results.Next() {
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
	ctx, cancle := context.WithTimeout(context.Background(), 15*time.Second) //if the query is going to take longer than 15 sec then its going to cancle and return
	defer cancle()
	if product.ProductID == 0 {
		return errors.New("product has invalid ID")
	}
	_, err := database.DbConn.ExecContext(ctx, `UPDATE products SET 
	manufacturer=?, 
	sku=?, 
	upc=?, 
	pricePerUnit=?,
	quantityOnHand=?,
	productName=?
	WHERE productId=?`,
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

func searchForProductData(productFilter ProductReportFilter) ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var queryArgs = make([]interface{}, 0)
	var queryBuilder strings.Builder
	queryBuilder.WriteString(`SELECT 
		productId, 
		LOWER(manufacturer), 
		LOWER(sku), 
		upc, 
		pricePerUnit, 
		quantityOnHand, 
		LOWER(productName) 
		FROM products WHERE `)
	if productFilter.NameFilter != "" {
		queryBuilder.WriteString(`productName LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(productFilter.NameFilter)+"%")
	}
	if productFilter.ManufacturerFilter != "" {
		if len(queryArgs) > 0 {
			queryBuilder.WriteString(" AND ")
		}
		queryBuilder.WriteString(`manufacturer LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(productFilter.ManufacturerFilter)+"%")
	}
	if productFilter.SKUFilter != "" {
		if len(queryArgs) > 0 {
			queryBuilder.WriteString(" AND ")
		}
		queryBuilder.WriteString(`sku LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(productFilter.SKUFilter)+"%")
	}

	results, err := database.DbConn.QueryContext(ctx, queryBuilder.String(), queryArgs...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	products := make([]Product, 0)
	for results.Next() {
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

func insertProduct(product Product) (int, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 15*time.Second) //if the query is going to take longer than 15 sec then its going to cancle and return
	defer cancle()
	result, err := database.DbConn.ExecContext(ctx, `INSERT INTO products
		(manufacturer,
		sku,
		upc, 
		pricePerUnit,
		quantityOnHand,
		productName) VALUES (?,?,?,?,?,?)`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName) //inserting into the database using Exec
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return 0, nil
	}
	return int(insertID), nil
}
