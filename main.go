package main

/*
Code By Nurul Hidayat
email : crmspy@gmail.com
*/

import (
	"github.com/crmspy/ssinventory/library/autoload"
	."github.com/crmspy/ssinventory/controller"
	"github.com/gin-gonic/gin"
)

func init() {
	autoload.Run()
	MigrateDb()
}

// main
func main() {
	router := gin.Default()
    
	v1 := router.Group("/api/v1/")
	{
		//product
		v1.GET("product", FetchAllMproduct)
		v1.POST("product", CreateMproduct)
		v1.PUT("product/:id", UpdateMproduct)
        v1.DELETE("product/:id", DeleteMproduct)
        
		//get all order data
		v1.GET("order", FetchAllTorder)
        v1.POST("order", CreateTorder)
        
		//inventory
        v1.GET("inventory", FetchAllMinventory)
        v1.POST("inventory", CreateMinventory)
        v1.PUT("inventory/:id", UpdateMinventory)
        v1.DELETE("inventory/:id", DeleteMinventory)
        v1.POST("inventory/migration", MigrationDataInventory)
        v1.POST("inventory/inout", Inout)

        //report
        v1.GET("inventory/availablestock", AvailableStock)
        v1.GET("inventory/goodreceipt", GoodReceipt)
        v1.GET("inventory/goodshipment", GoodShipment)
        v1.GET("inventory/valueofproduct", ValueofProduct)
        v1.POST("inventory/salesorder", SalesOrder)
	}
	router.Run()
}
