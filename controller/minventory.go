package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
    "github.com/crmspy/ssinventory/library/conn"
    "github.com/crmspy/ssinventory/library/inventory"
    "time"
    "fmt"
    "strconv"
)

type (
    //inventory location
	mInventory struct {
		M_inventory_id	        string		`gorm:"type:varchar(64);PRIMARY_KEY"`
		Name			        string		`gorm:"type:varchar(255);"`
	}

    //information all product at inventory
	mInventoryLine struct {
		M_inventory_line_id		int			`gorm:"AUTO_INCREMENT"`
		M_inventory_id			string		`gorm:"type:varchar(64);"`
		M_product_id			string		`gorm:"type:varchar(64);"`
        Qty_count               int         `gorm:"type:int"`
        Last_update             time.Time   `gorm:"type:datetime"`
	}

    //Inout product at inventory
    tInout struct {
        T_inout_id		        int			`gorm:"AUTO_INCREMENT"`
        Inout_type		        string		`gorm:"type:enum(3);"`
		M_inventory_id			string		`gorm:"type:varchar(64);"`
        M_product_id			string		`gorm:"type:varchar(64);"`
        T_order_line_id         int         `gorm:"type:int"`
        Inout_qty               int         `gorm:"type:int"`
        Inout_date              time.Time   `gorm:"type:datetime"`
    }

    // transformedmInventory represents a formatted inventory location
	transformedmInventory struct {
		M_inventory_id	string		`json:"inventory_id"`
		Name			string		`json:"name"`
    }
    
)

// Create Inventory location
func CreateMinventory(c *gin.Context) {
	model_mInventory := mInventory{
        M_inventory_id: c.PostForm("m_inventory_id"),
        Name: c.PostForm("name"),
    }
	conn.Db.Save(&model_mInventory)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusOK, "message": "inventory location created successfully!", "resourceId": model_mInventory.M_inventory_id})
}

// fetch all inventory location
func FetchAllMinventory(c *gin.Context) {
    var modelmInventory []mInventory
    var _modelmInventory []transformedmInventory
    conn.Db.Find(&modelmInventory)
    if len(modelmInventory) <= 0 {
        c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Inventory location Not found!"})
        return
    }
    //transforms the Inventory for building a good response
    for _, item := range modelmInventory {
        _modelmInventory = append(_modelmInventory, transformedmInventory{
            M_inventory_id: item.M_inventory_id,
            Name: item.Name,
        })
    }
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _modelmInventory})
}


// remove a inventory location
func DeleteMinventory(c *gin.Context) {
    var modelmInventory mInventory
    m_inventory_id := c.Param("id")

    conn.Db.Where("m_inventory_id = ?", m_inventory_id).First(&modelmInventory)
    if modelmInventory.M_inventory_id == "" {
        c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Inventory location found!"})
        return
        }
    conn.Db.Delete(&modelmInventory)
        c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Inventory location deleted successfully!"})
}


/*
the main function of the inventory application
this API Serve to update the stock of goods out / in from the warehouse
please use IN / OUT for inout type
IN = Put the product into the warehouse
OUT = Eject product from the warehouse
*/
func Inout(c *gin.Context){
    qty,_ := strconv.Atoi(c.PostForm("qty"))
    t_order_line_id,_ := strconv.Atoi(c.PostForm("t_order_line_id"))
    //test inventory manual input
    param := inventory.InoutParam {
        M_inventory_id: c.PostForm("m_inventory_id"),
        Inout_qty: qty,
        T_order_line_id: t_order_line_id,
    }
    if e := inventory.DoInout(param); e != nil{
        fmt.Println(e)
        c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprint(e)})
    }else{
        c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "inventory '"+param.M_inventory_id+"' was updated successfully"})
    }
}
