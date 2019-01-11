package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
    "github.com/crmspy/ssinventory/library/conn"
    ."github.com/crmspy/ssinventory/modules/inventory/models"
    "errors"
)

// Create Product
func CreateMproduct(c *gin.Context) {
    m_product_id := c.PostForm("m_product_id")
    if err := IssetProduct(m_product_id); err == nil{
        c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "data": "Product "+m_product_id+" already exist"})
        return
    }
	modelProduct := Mproduct{
        M_product_id: m_product_id,
        Name: c.PostForm("name"),
    }
	conn.Db.Save(&modelProduct)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Product created successfully!", "resourceId": modelProduct.M_product_id})
}

// fetch all product
func FetchAllMproduct(c *gin.Context) {
    var modelProduct []Mproduct
    var _modelProduct []TransformedMproduct
    conn.Db.Find(&modelProduct)
    if len(modelProduct) <= 0 {
        c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No product found!"})
        return
    }
    //transforms the product for building a good response
    for _, item := range modelProduct {
        _modelProduct = append(_modelProduct, TransformedMproduct{
        M_product_id: item.M_product_id,
        Name: item.Name,
    })
    }
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _modelProduct})
}


func IssetProduct(id string) (error){
    var M_product_id string
    conn.Db.Table("m_product").Where("m_product_id = ?", id).Select("m_product_id").Row().Scan(&M_product_id)
    if M_product_id=="" {
        return errors.New("Cannot find product in out database")
    }
    return nil
}

// update a product
func UpdateMproduct(c *gin.Context) {
    var modelMproduct Mproduct
    m_product_id := c.Param("id")

    conn.Db.First(&modelMproduct, "m_product_id = ?", m_product_id)

    if modelMproduct.M_product_id == "" {
        c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No product found!"})
        return
    }

    conn.Db.Model(&modelMproduct).Update("name", c.PostForm("name"))
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Product updated successfully!"})
}

//delete product
func DeleteMproduct(c *gin.Context) {
    var modelMproduct Mproduct
    m_product_id := c.Param("id")

    conn.Db.Where("m_product_id = ?", m_product_id).First(&modelMproduct)
    if modelMproduct.M_product_id == "" {
        c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No product found!"})
        return
        }
    conn.Db.Delete(&modelMproduct)
        c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Product deleted successfully!"})
}

