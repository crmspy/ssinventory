package main

/*
Code By Nurul Hidayat
email : crmspy@gmail.com
*/

import (
	"github.com/crmspy/ssinventory/library/autoload"
	."github.com/crmspy/ssinventory/controller"
	."github.com/crmspy/ssinventory/library/auth"
	"github.com/gin-gonic/gin"
)

func init() {
	autoload.Run()
	MigrateDb()
}

// main
func main() {
	router := gin.Default()
	auth := router.Group("/api/user")
	{
		auth.POST("/login", GetKey)
        auth.GET("/getProfile", Auth,GetProfile)
        /*
        v1.Use(Auth) call middleware
        auth.GET("/logout", Auth,GetTodo)
		auth.GET("/session", Auth,GetTodo) //get active session
        auth.DELETE("/session/:id", Auth,GetTodo) //delete session
        */
	}
	v0 := router.Group("/api/v1/")
	{
		//product
		v0.GET("product", FetchAllMproduct)
		v0.POST("product", CreateMproduct)

		//get all order data
		v0.GET("order", FetchAllTorder)
        v0.POST("order", CreateTorder)
        
		//inventory
        v0.GET("inventory", FetchAllMinventory)
		v0.POST("inventory", CreateMinventory)
        v0.DELETE("inventory/:id", DeleteMinventory)
        v0.POST("inventory/migration", MigrationDataInventory)
        v0.POST("inventory/inout", Inout)

        //report
        v0.GET("inventory/availablestock", AvailableStock)
        v0.GET("inventory/goodreceipt", GoodReceipt)
        v0.GET("inventory/goodshipment", GoodShipment)
        v0.GET("inventory/valueofproduct", ValueofProduct)
        v0.POST("inventory/salesorder", SalesOrder)
	}
	router.Run()
}
