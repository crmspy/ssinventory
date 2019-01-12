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
)

type (
	tOrder struct {
		//gorm.Model
		T_order_id		string		`gorm:"type:varchar(64);PRIMARY_KEY"`
		Order_type		string		`gorm:"type:varchar(1);"`
        Order_amount	float64		`gorm:"type:float(16,2);"`
        Description     string
        Order_status	string		`gorm:"type:varchar(1);"`
        Order_date      time.Time   `gorm:"type:datetime"`
    }

    tOrderLine struct {
		//gorm.Model
		T_order_line_id		        int		    `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
		T_order_id  		        string		`gorm:"type:varchar(64);"`
        M_product_id	            string		`gorm:"type:varchar(64);"`
        Orderline_qty               int         `gorm:"type:int"`
        Orderline_price             float64		`gorm:"type:float(16,2);"`
        Orderline_total_amount	    float64		`gorm:"type:float(16,2);"`
        Orderline_outstanding       int         `gorm:"type:int"`
        Orderline_received          int         `gorm:"type:int"`
    }
    
    // transformed Order represents a formatted product
	transformedtOrder struct {
		T_order_id		string		`json:"transaction_id"`
		Order_type		string		`json:"order_type"`
        Order_amount	float64		`json:"amount"`
        Description     string      `json:"description"`
        Order_status    string      `json:"order_status"`
        Order_date      time.Time   `json:"order_date"`
    }
    
    transformedtOrderLine struct {
		T_order_line_id		        string		`json:"t_order_line_id"`
		T_order_id  		        string		`json:"t_order_id"`
        M_product_id	            string		`json:"m_product_id"`
        Orderline_qty               int         `json:"qty"`
        Orderline_price             float64		`json:"price"`      
        Orderline_total_amount	    float64		`json:"total_amount"`
        Orderline_outstanding       int         `json:"outstanding"`
        Orderline_received          int         `json:"received"`
    }
)


// fetch all order
func FetchAllTorder(c *gin.Context) {
    var modeltOrder []tOrder
    var _modeltOrder []transformedtOrder
    page,_  := strconv.Atoi(c.DefaultQuery("page", "1"))
    offset,limit := cgx.Calcpage(page)
    conn.Db.Offset(offset).Limit(limit).Find(&modeltOrder)
    if len(modeltOrder) <= 0 {
        c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Order found!"})
        return
    }
    //transforms the Order for building a good response
    for _, item := range modeltOrder {
        _modeltOrder = append(_modeltOrder, transformedtOrder{
        T_order_id: item.T_order_id,
        Order_type: item.Order_type,
        Order_amount: item.Order_amount,
        Description: item.Description,
        Order_status: item.Order_status,
        Order_date: item.Order_date,
    })
    }
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _modeltOrder})
}

/*create order
this is sampe parameter
`{"detail":[{"m_product_id":"acessToken","qty":1000,"price":1000.00,"total_amount":100000}]}`
*/
func CreateTorder(c *gin.Context) {

    lineData := c.PostForm("detail")
	lineJSON := make(map[string][]transformedtOrderLine)
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
            `select t_order_id from t_orders where t_order_id = ? limit 1`,t_order_id).Row().Scan(&t_order_id_exist)
        if (t_order_id_exist != ""){
            c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "data": "t_order_id already exist, leave it blank if you wanna automatically generate it"})
            return
        }
    }

    modeltOrder := tOrder{
        T_order_id: t_order_id,
        Order_type: order_type,
        Description: c.PostForm("description"),
        Order_status: c.PostForm("order_status"),
        Order_date  : time.Now(),
    }

    if err := tx.Save(&modeltOrder).Error; err != nil {
       tx.Rollback()
       c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "data": "Failed create order"})
       return
    }

    var total_amount float64
    for _,value := range lineJSON["detail"]{
        amount := float64 (float64(value.Orderline_qty) *  value.Orderline_price)
        if err := IssetProduct(value.M_product_id); err != nil{
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "data": "Product "+value.M_product_id+" not found"})
            return
        }

        modeltOrderLine := tOrderLine{
            T_order_id              : modeltOrder.T_order_id,
            M_product_id            : value.M_product_id,
            Orderline_qty           : value.Orderline_qty,
            Orderline_price         : value.Orderline_price,
            Orderline_total_amount  : amount,
            Orderline_outstanding   : value.Orderline_qty,
        }
        
        //insert product order
        if err := tx.Save(&modeltOrderLine).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "data": "Failed insert line order"})
        }
        total_amount += amount
    }

    //update total order
    modeltOrder.Order_amount = total_amount
    tx.Save(&modeltOrder)
    tx.Commit()
    
    c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "order data successfully!", "resourceId": modeltOrder.T_order_id})

}
