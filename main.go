package main

import (
	"net/http"

	"github.com/indikamaligaspe/inventoryservice/database"

	_ "github.com/go-sql-driver/mysql"
	"github.com/indikamaligaspe/inventoryservice/product"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
