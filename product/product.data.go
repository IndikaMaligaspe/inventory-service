package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/indikamaligaspe/inventoryservice/database"
)

var productMap = struct {
	sync.RWMutex
	m map[int]Product
}{m: make(map[int]Product)}

func init() {
	fmt.Println("loading products.......")
	prodMap, err := loadProductMap()
	productMap.m = prodMap
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d products loaded... \n", len(productMap.m))
}

func loadProductMap() (map[int]Product, error) {
	fileName := "products.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}
	file, _ := ioutil.ReadFile(fileName)
	productList := make([]Product, 0)
	err = json.Unmarshal([]byte(file), &productList)
	if err != nil {
		log.Fatal(err)
	}
	prodMap := make(map[int]Product)
	for i := 0; i < len(productList); i++ {
		prodMap[productList[i].ProductID] = productList[i]
	}
	return prodMap, nil
}

func getProduct(productID int) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
	productId, 
	manufacturer, 
	sku, 
	upc, 
	pricePerUnit, 
	quantityOnHand, 
	productName 
	FROM products
	WHERE productId = ?
	`, productID)
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
		log.Fatal(err)
		return nil, err
	}
	return product, nil
}

func getProductList() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result, err := database.DbConn.QueryContext(ctx, `SELECT 
	productId, 
	manufacturer, 
	sku, 
	upc, 
	pricePerUnit, 
	quantityOnHand, 
	productName 
	FROM products`)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer result.Close()
	products := make([]Product, 0)
	for result.Next() {
		var product Product
		result.Scan(&product.ProductID,
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
	_, err := database.DbConn.Exec(`UPDATE products SET 
		manufacturer=?, 
		sku=?, 
		upc=?, 
		pricePerUnit=CAST(? AS DECIMAL(13,2)), 
		quantityOnHand=?, 
		productName=?
		WHERE productId=?`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName,
		product.ProductID)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func addProduct(product Product) (int, error) {
	result, err := database.DbConn.Exec(`INSERT INTO products  
		(manufacturer, 
		sku, 
		upc, 
		pricePerUnit, 
		quantityOnHand, 
		productName) VALUES (?, ?, ?, ?, ?, ?)`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName)
	if err != nil {
		log.Fatal(err.Error())
		return 0, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err.Error())
		return 0, err
	}
	return int(insertID), nil
}

func removeProduct(productID int) error {
	_, err := database.DbConn.Exec(`DELETE FROM products WHERE productId=?`, productID)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
