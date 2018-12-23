package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
    "github.com/crmspy/ssinventory/library/conn"
    "errors"
)

type (
	mProduct struct {
		M_product_id	string		`gorm:"type:varchar(64);PRIMARY_KEY"`
		Name			string		`gorm:"type:varchar(255);"`
	}
    // transformedmProductrepresents a formatted product
	transformedmProduct struct {
		M_product_id	string		`json:"sku"`
		Name			string		`json:"product_name"`
	}
)

// Create Product
func CreateMproduct(c *gin.Context) {
    m_product_id := c.PostForm("m_product_id")
    if err := IssetProduct(m_product_id); err == nil{
        c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "data": "Product "+m_product_id+" already exist"})
        return
    }
	model_mProduct := mProduct{
        M_product_id: m_product_id,
        Name: c.PostForm("name"),
    }
	conn.Db.Save(&model_mProduct)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Product created successfully!", "resourceId": model_mProduct.M_product_id})
}

// fetch all product
func FetchAllMproduct(c *gin.Context) {
    var modelmProduct []mProduct
    var _modelmProduct []transformedmProduct
    conn.Db.Find(&modelmProduct)
    if len(modelmProduct) <= 0 {
        c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
        return
    }
    //transforms the product for building a good response
    for _, item := range modelmProduct {
        _modelmProduct = append(_modelmProduct, transformedmProduct{
        M_product_id: item.M_product_id,
        Name: item.Name,
    })
    }
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _modelmProduct})
}

func IssetProduct(id string) (error){
    var M_product_id string
    conn.Db.Table("m_products").Where("m_product_id = ?", id).Select("m_product_id").Row().Scan(&M_product_id)
    if M_product_id=="" {
        return errors.New("Cannot find product in out database")
    }
    return nil
}

// update a product
func UpdateMproduct(c *gin.Context) {
    var modelmProduct mProduct
    m_product_id := c.Param("id")

    conn.Db.First(&modelmProduct, "m_product_id = ?", m_product_id)

    if modelmProduct.M_product_id == "" {
        c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No product found!"})
        return
    }

    conn.Db.Model(&modelmProduct).Update("name", c.PostForm("name"))
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Product updated successfully!"})
}

//delete product
func DeleteMproduct(c *gin.Context) {
    var modelmProduct mProduct
    m_product_id := c.Param("id")

    conn.Db.Where("m_product_id = ?", m_product_id).First(&modelmProduct)
    if modelmProduct.M_product_id == "" {
        c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No product found!"})
        return
        }
    conn.Db.Delete(&modelmProduct)
        c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Product deleted successfully!"})
}
