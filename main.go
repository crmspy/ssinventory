package main

/*
Code By Nurul Hidayat
email : crmspy@gmail.com
*/

import (
	"github.com/crmspy/ssinventory/library/autoload"
    inventory "github.com/crmspy/ssinventory/modules/inventory/controller"
	order "github.com/crmspy/ssinventory/modules/order/controller"
	"github.com/gin-gonic/gin"
)

func init() {
	autoload.Run()
}

// main
func main() {
    router := gin.Default()
    v1 := router.Group("/api/v1/")
	{
        
		//product
		v1.GET("product", inventory.FetchAllMproduct)
		v1.POST("product", inventory.CreateMproduct)
		v1.PUT("product/:id", inventory.UpdateMproduct)
        v1.DELETE("product/:id", inventory.DeleteMproduct)

        //get all order data
		v1.GET("order", order.FetchAllTorder)
        v1.POST("order", order.CreateTorder)
        v1.POST("order/status", order.UpdateStatus)

        //inventory
        v1.GET("inventory", inventory.FetchAllMinventory)
        v1.POST("inventory", inventory.CreateMinventory)
        v1.PUT("inventory/:id", inventory.UpdateMinventory)
        v1.DELETE("inventory/:id", inventory.DeleteMinventory)

        v1.POST("inventory/inout", inventory.Inout)
        v1.POST("inventory/uploadstock", inventory.UploadStock)

        //report
        v1.GET("inventory/availablestock", inventory.AvailableStock)
        v1.GET("inventory/goodreceipt", inventory.GoodReceipt)
        v1.GET("inventory/goodshipment", inventory.GoodShipment)
        v1.GET("inventory/valueofproduct", inventory.ValueofProduct)
        v1.POST("inventory/salesorder", inventory.SalesOrder)
    }
	router.Run()
}
