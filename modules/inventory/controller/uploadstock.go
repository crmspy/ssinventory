package controller
// create tabel in database if not exist
import (
	"github.com/gin-gonic/gin"
    "github.com/crmspy/ssinventory/library/conn"
    "encoding/csv"
    "os"
    "time"
    "net/http"
    log "github.com/Sirupsen/logrus"
    "strconv"
    order "github.com/crmspy/ssinventory/modules/order/models"
    ."github.com/crmspy/ssinventory/modules/inventory/models"
)

/* Migration data inventory
this migration process takes several steps:
1. creating order dummy so system can get information about product, purchase price & last stock
2. move order data to inventory
*/
func UploadStock(c *gin.Context) {
    mycsv := readFile(c)

    modelTorder := order.Torder{
        T_order_id: "IMPORTCSV"+time.Now().Format("200601021504"),
        Order_type: "P",
        Description: "Import From Csv",
        Order_status: "C",
        Order_date  : time.Now(),
    }

    if err := conn.Db.Save(&modelTorder).Error; err != nil {
       c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "data": "Failed import stock"})
       return
    }


    /*
        column[0] = Product Code
        column[1] = Product Name
        column[2] = QTY
        column[3] = Price
        column[4] = Inventory
    */
    var total_amount float64
    for key, column := range mycsv {
        if key > 0 {
                time.Sleep(50000)
                //insert product
                modelMproduct := Mproduct{
                    M_product_id: column[0],
                    Name: column[1],
                }

                conn.Db.Save(&modelMproduct)
                
                //add product to list order
                qty,_ := strconv.ParseInt(column[2],0,64)
                price, _ := strconv.ParseFloat(column[3], 64)
                
                amount := float64 (float64(qty) *  price)
                modelTorderLine := order.TorderLine{
                    T_order_id              : modelTorder.T_order_id,
                    M_product_id            : column[0],
                    Orderline_qty           : int (qty),
                    Orderline_price         : price,
                    Orderline_total_amount  : amount,
                    Orderline_outstanding   : int (qty),
                }
                conn.Db.Save(&modelTorderLine)
                
                var M_inventory_id string
                if column[4]==""{
                    M_inventory_id = "General"
                }else{
                    M_inventory_id = column[4]
                }
                //updating stock in inventory
                param := InoutParam {
                    M_inventory_id  : M_inventory_id,
                    Inout_qty       : int (qty),
                    T_order_line_id : int (modelTorderLine.T_order_line_id),
                    Description     : "Import From CSV",
                }
                DoInout(param)
                total_amount += amount
        }
    }

    //update total order
    modelTorder.Order_amount = total_amount
    conn.Db.Save(&modelTorder)

    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Inventory data successfully uploaded"})
}

func readFile(c *gin.Context) (mycsv [][]string){
    file, err := c.FormFile("file")
    if err != nil {
        log.Fatal(err)
    }
    log.Println(file.Filename)

    err = c.SaveUploadedFile(file, "tmp/"+file.Filename)
    if err != nil {
        log.Fatal(err)
    }

    // Open CSV file
    f, err := os.Open("tmp/"+file.Filename)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    // Read File into a Variable
    lines, err := csv.NewReader(f).ReadAll()
    if err != nil {
        panic(err)
    }
    return lines

}
