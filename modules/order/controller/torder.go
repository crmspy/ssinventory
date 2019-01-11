package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
    "github.com/crmspy/ssinventory/library/conn"
    "github.com/crmspy/ssinventory/library/cgx"
    "encoding/json"
    "time"
    "math/rand"
    "fmt"
    "strconv"
    ."github.com/crmspy/ssinventory/modules/order/models"
    inventory "github.com/crmspy/ssinventory/modules/inventory/controller"
)

// fetch all order
func FetchAllTorder(c *gin.Context) {
    var modelTorder []Torder
    var _modelTorder []TransformedTorder
    page,_  := strconv.Atoi(c.DefaultQuery("page", "1"))
    offset,limit := cgx.Calcpage(page)

    conn.Db.Offset(offset).Limit(limit).Find(&modelTorder)
    if len(modelTorder) <= 0 {
        c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Order found!"})
        return
    }
    //transforms the Order for building a good response
    for _, item := range modelTorder {
        _modelTorder = append(_modelTorder, TransformedTorder{
        T_order_id: item.T_order_id,
        Order_type: item.Order_type,
        Order_amount: item.Order_amount,
        Description: item.Description,
        Order_status: item.Order_status,
        Order_date: item.Order_date,
    })
    }
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _modelTorder})
}

/*create order
this is sampe parameter
`{"detail":[{"m_product_id":"acessToken","qty":1000,"price":1000.00,"total_amount":100000}]}`
*/
func CreateTorder(c *gin.Context) {

    lineData := c.PostForm("detail")
	lineJSON := make(map[string][]TransformedTorderLine)
	err := json.Unmarshal([]byte(lineData), &lineJSON)
	if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "data": "Format data json not valid"})
        return
    }
    
    //this transaction method, if failed data will rollback
    tx := conn.Db.Begin()
    
    var t_order_id string = c.PostForm("t_order_id");
    var order_type string = c.PostForm("order_type");

    if order_type=="" {
        c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "data": "order_type must be set, you can use PO = Purchase Order, SO = Sales Order"})
        return
    }

    if t_order_id ==""{
        rand.Seed(time.Now().Unix())
        random_number :=  fmt.Sprintf("%v",rand.Intn(1000 - 13) + 13)
        if order_type == "S"{
            t_order_id = "SO_SS"+random_number+""+time.Now().Format("200601021504");
        }else{
            t_order_id = "PO_SS"+random_number+""+time.Now().Format("200601021504");
        }
    }else{
        var t_order_id_exist string

        //check t_order_id already exist
        conn.Db.Raw(
            `select t_order_id from t_order where t_order_id = ? limit 1`,t_order_id).Row().Scan(&t_order_id_exist)
        if (t_order_id_exist != ""){
            c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "data": "t_order_id already exist, leave it blank if you wanna automatically generate it"})
            return
        }
    }

    modelTorder := Torder{
        T_order_id: t_order_id,
        Order_type: order_type,
        Description: c.PostForm("description"),
        Order_status: c.PostForm("order_status"),
        Order_date  : time.Now(),
    }

    if err := tx.Save(&modelTorder).Error; err != nil {
       tx.Rollback()
       c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "data": "Failed create order"})
       return
    }

    var total_amount float64
    for _,value := range lineJSON["detail"]{
        amount := float64 (float64(value.Orderline_qty) *  value.Orderline_price)
        if err := inventory.IssetProduct(value.M_product_id); err != nil{
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "data": "Product "+value.M_product_id+" not found"})
            return
        }

        modelTorderLine := TorderLine{
            T_order_id              : modelTorder.T_order_id,
            M_product_id            : value.M_product_id,
            Orderline_qty           : value.Orderline_qty,
            Orderline_price         : value.Orderline_price,
            Orderline_total_amount  : amount,
            Orderline_outstanding   : value.Orderline_qty,
        }
        
        //insert product order
        if err := tx.Save(&modelTorderLine).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "data": "Failed insert line order"})
        }
        total_amount += amount
    }

    //update total order
    modelTorder.Order_amount = total_amount
    tx.Save(&modelTorder)
    tx.Commit()
    
    c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "order data successfully!", "resourceId": modelTorder.T_order_id})

}
