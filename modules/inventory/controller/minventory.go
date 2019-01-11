package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
    "github.com/crmspy/ssinventory/library/conn"
    ."github.com/crmspy/ssinventory/modules/inventory/models"
)


// Create Inventory location
func CreateMinventory(c *gin.Context) {
    var model_Minventory Minventory
    m_inventory_id := c.PostForm("m_inventory_id")

    conn.Db.First(&model_Minventory, "m_inventory_id = ?", m_inventory_id)
    if model_Minventory.M_inventory_id != "" {
        c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Inventory location already exist!"})
        return
    }

	model_Minventory = Minventory{
        M_inventory_id: m_inventory_id,
        Name: c.PostForm("name"),
    }
	conn.Db.Save(&model_Minventory)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusOK, "message": "inventory location created successfully!", "resourceId": model_Minventory.M_inventory_id})
}

// fetch all inventory location
func FetchAllMinventory(c *gin.Context) {
    var modelMinventory []Minventory
    var _modelMinventory []TransformedMinventory
    conn.Db.Find(&modelMinventory)
    if len(modelMinventory) <= 0 {
        c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Inventory location Not found!"})
        return
    }
    //transforms the Inventory for building a good response
    for _, item := range modelMinventory {
        _modelMinventory = append(_modelMinventory, TransformedMinventory{
            M_inventory_id: item.M_inventory_id,
            Name: item.Name,
        })
    }
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _modelMinventory})
}


// update a inventory location
func UpdateMinventory(c *gin.Context) {
    var modelMinventory Minventory
    m_inventory_id := c.Param("id")

    conn.Db.First(&modelMinventory, "m_inventory_id = ?", m_inventory_id)

    if modelMinventory.M_inventory_id == "" {
        c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No inventory location found!"})
        return
    }

    conn.Db.Model(&modelMinventory).Update("name", c.PostForm("name"))
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Inventory Location updated successfully!"})
}

// remove a inventory location
func DeleteMinventory(c *gin.Context) {
    var modelMinventory Minventory
    m_inventory_id := c.Param("id")

    conn.Db.Where("m_inventory_id = ?", m_inventory_id).First(&modelMinventory)
    if modelMinventory.M_inventory_id == "" {
        c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Inventory location found!"})
        return
        }
    conn.Db.Delete(&modelMinventory)
        c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Inventory location deleted successfully!"})
}
