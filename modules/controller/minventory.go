package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
    "github.com/crmspy/ssinventory/library/conn"
    "github.com/crmspy/ssinventory/library/inventory"
    "time"
    "fmt"
    "strconv"
    "bytes"
    log "github.com/Sirupsen/logrus"
    "encoding/csv"
    
)

type (
    //inventory location
	mInventory struct {
		M_inventory_id	        string		`gorm:"type:varchar(64);PRIMARY_KEY"`
		Name			        string		`gorm:"type:varchar(255);"`
	}

    //information all product at inventory
	mInventoryLine struct {
		M_inventory_line_id		int			`gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
		M_inventory_id			string		`gorm:"type:varchar(64);"`
		M_product_id			string		`gorm:"type:varchar(64);"`
        Qty_count               int         `gorm:"type:int"`
        Last_update             time.Time   `gorm:"type:datetime"`
	}

    //Inout product at inventory
    tInout struct {
        T_inout_id		        int			`gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
        Inout_type		        string		`gorm:"type:varchar(3);"`
		M_inventory_id			string		`gorm:"type:varchar(64);"`
        M_product_id			string		`gorm:"type:varchar(64);"`
        T_order_line_id         int         `gorm:"type:int"`
        Inout_qty               int         `gorm:"type:int"`
        Inout_date              time.Time   `gorm:"type:datetime"`
        Description             string      `gorm:"type:varchar(255)"`
    }

    // transformedmInventory represents a formatted inventory location
	transformedmInventory struct {
		M_inventory_id	string		`json:"m_inventory_id"`
		Name			string		`json:"name"`
    }
    
)

// Create Inventory location
func CreateMinventory(c *gin.Context) {
    var model_mInventory mInventory
    m_inventory_id := c.PostForm("m_inventory_id")

    conn.Db.First(&model_mInventory, "m_inventory_id = ?", m_inventory_id)
    if model_mInventory.M_inventory_id != "" {
        c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Inventory location already exist!"})
        return
    }

	model_mInventory = mInventory{
        M_inventory_id: m_inventory_id,
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


// update a inventory location
func UpdateMinventory(c *gin.Context) {
    var modelmInventory mInventory
    m_inventory_id := c.Param("id")

    conn.Db.First(&modelmInventory, "m_inventory_id = ?", m_inventory_id)

    if modelmInventory.M_inventory_id == "" {
        c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No inventory location found!"})
        return
    }

    conn.Db.Model(&modelmInventory).Update("name", c.PostForm("name"))
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Inventory Location updated successfully!"})
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
        Description: c.PostForm("description"),
    }
    if e := inventory.DoInout(param); e != nil{
        c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprint(e)})
    }else{
        c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "inventory '"+param.M_inventory_id+"' was updated successfully"})
    }
}


//download available product in inventory
func AvailableStock(c *gin.Context){

    b := &bytes.Buffer{}
    w := csv.NewWriter(b)

    //header csv
    if err := w.Write([]string{"Inventory Location","SKU", "Nama Item","Jumlah Barang","Last Update"}); err != nil {
        log.Fatalln("error writing record to csv:", err)
    }

    // get available stock
    rows, _ := conn.Db.Raw(`
    select 
        m_products.m_product_id,
        m_products.name,
        m_inventory_lines.qty_count,
        m_inventory_lines.m_inventory_id,
        m_inventory_lines.last_update
    from 
        m_inventory_lines 
        left join m_products on (m_products.m_product_id = m_inventory_lines.m_product_id) 
        order by m_inventory_lines.m_product_id asc
    `).Rows() // (*sql.Rows, error)
    defer rows.Close()
    for rows.Next() {
        var m_product_id,product_name,qty_count,m_inventory_id string
        var last_update time.Time
        rows.Scan(&m_product_id,&product_name,&qty_count,&m_inventory_id,&last_update)

        var record []string
        record = append(record, m_inventory_id)
        record = append(record, m_product_id)
        record = append(record, product_name)
        record = append(record, qty_count)
        record = append(record, last_update.Format("2006-01-02 15:04:05"))
        if err := w.Write(record); err != nil {
            log.Fatalln("error writing record to csv:", err)
        }
    }
    w.Flush()

    if err := w.Error(); err != nil {
        log.Fatal(err)
    }
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Disposition", "attachment; filename=AvailableStock.csv")
    c.Data(http.StatusOK, "text/csv", b.Bytes())
}

//download good receipt
func GoodReceipt(c *gin.Context){

    b := &bytes.Buffer{}
    w := csv.NewWriter(b)
    //header csv
    if err := w.Write([]string{"Waktu","SKU", "Nama Barang","Jumlah Pesanan","Jumlah Diterima","Harga Beli","Total","No Kwitansi","Catatan"}); err != nil {
        log.Fatalln("error writing record to csv:", err)
    }

    // get available stock
    rows, _ := conn.Db.Raw(`
        select
            t_inouts.inout_date,
            m_products.m_product_id,
            m_products.name,
            t_order_lines.orderline_qty,
            t_order_lines.orderline_received,
            t_order_lines.orderline_price,
            t_order_lines.orderline_total_amount,
            t_orders.t_order_id,
            t_inouts.description
        from 
            t_inouts
            left join t_order_lines on (t_order_lines.t_order_line_id = t_inouts.t_order_line_id)
            left join m_products on (m_products.m_product_id = t_order_lines.m_product_id)
            left join t_orders on (t_orders.t_order_id = t_order_lines.t_order_id)
            where t_inouts.inout_type = 'IN'
            group by t_inouts.t_order_line_id
            order by t_inouts.inout_date asc
        `).Rows() // (*sql.Rows, error)
    defer rows.Close()
    for rows.Next() {
        var m_product_id,name,orderline_qty,orderline_received,orderline_price,orderline_total_amount,t_order_id,description string
        var inout_date time.Time
        rows.Scan(&inout_date,&m_product_id,&name,&orderline_qty,&orderline_received,&orderline_price,&orderline_total_amount,&t_order_id,&description)
        var record []string
        record = append(record, inout_date.Format("2006-01-02 15:04:05"))
        record = append(record, m_product_id)
        record = append(record, name)
        record = append(record, orderline_qty)
        record = append(record, orderline_received)
        record = append(record, orderline_price)
        record = append(record, orderline_total_amount)
        record = append(record, t_order_id)
        record = append(record, description)
        if err := w.Write(record); err != nil {
            log.Fatalln("error writing record to csv:", err)
        }
    } 
    w.Flush()

    if err := w.Error(); err != nil {
        log.Fatal(err)
    }
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Disposition", "attachment; filename=GoodReceipt.csv")
    c.Data(http.StatusOK, "text/csv", b.Bytes())
}

//download good shipment
func GoodShipment(c *gin.Context){

    b := &bytes.Buffer{}
    w := csv.NewWriter(b)
    //header csv
    if err := w.Write([]string{"Waktu","SKU", "Nama Barang","Jumlah Pesanan","Jumlah Dikirim","Harga Jual","Total","No Pesanan","Catatan"}); err != nil {
        log.Fatalln("error writing record to csv:", err)
    }

    // get available stock
    rows, _ := conn.Db.Raw(`
        select
            t_inouts.inout_date,
            m_products.m_product_id,
            m_products.name,
            t_order_lines.orderline_qty,
            t_order_lines.orderline_received,
            t_order_lines.orderline_price,
            t_order_lines.orderline_total_amount,
            t_orders.t_order_id,
            t_inouts.description
        from 
            t_inouts
            left join t_order_lines on (t_order_lines.t_order_line_id = t_inouts.t_order_line_id)
            left join m_products on (m_products.m_product_id = t_order_lines.m_product_id)
            left join t_orders on (t_orders.t_order_id = t_order_lines.t_order_id)
            where t_inouts.inout_type = 'OUT'
            group by t_inouts.t_order_line_id
            order by t_inouts.inout_date asc
        `).Rows() // (*sql.Rows, error)
    defer rows.Close()
    for rows.Next() {
        var m_product_id,name,orderline_qty,orderline_received,orderline_price,orderline_total_amount,t_order_id,description string
        var inout_date time.Time
        rows.Scan(&inout_date,&m_product_id,&name,&orderline_qty,&orderline_received,&orderline_price,&orderline_total_amount,&t_order_id,&description)
        var record []string
        record = append(record, inout_date.Format("2006-01-02 15:04:05"))
        record = append(record, m_product_id)
        record = append(record, name)
        record = append(record, orderline_qty)
        record = append(record, orderline_received)
        record = append(record, orderline_price)
        record = append(record, orderline_total_amount)
        record = append(record, t_order_id)
        record = append(record, description)
        if err := w.Write(record); err != nil {
            log.Fatalln("error writing record to csv:", err)
        }
    }
    w.Flush()

    if err := w.Error(); err != nil {
        log.Fatal(err)
    }
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Disposition", "attachment; filename=GoodShipment.csv")
    c.Data(http.StatusOK, "text/csv", b.Bytes())
}

//download report value of product
func ValueofProduct(c *gin.Context){

    b := &bytes.Buffer{}
    w := csv.NewWriter(b)

    //query get data transaction product
    sqlquery := `
    select 
        m_products.m_product_id,
        m_products.name,
        sum(orderline_qty) as orderline_qty,
        avg(orderline_price) as average_price,
        sum(orderline_total_amount) as orderline_total_amount
    from
        t_order_lines
        left join m_products on (m_products.m_product_id = t_order_lines.m_product_id)
        left join t_orders on (t_orders.t_order_id = t_order_lines.t_order_id)
        where t_orders.order_type = 'P'
        GROUP by t_order_lines.m_product_id
        order by t_order_lines.m_product_id asc
    `

    //prepare header data
    timenow := time.Now().Format("02-01-2006")
    var total_sku,total_product,total_amount_product string

    //subselect query to get summary data
    conn.Db.Raw(
        `select 
            count(*) as total_sku,
            sum(orderline_qty) as total_product,
            sum(orderline_total_amount) as total_amount_product 
        from (
            `+sqlquery+`
        ) as TB`).Row().Scan(&total_sku,&total_product,&total_amount_product)

    //header csv
    if err := w.Write([]string{"Laporan Nilai Barang"}); err != nil {
            log.Fatalln("error writing record to csv:", err)
    }
    w.Write([]string{"",});
    w.Write([]string{"Tanggal Cetak", timenow});
    w.Write([]string{"Jumlah SKU", total_sku});
    w.Write([]string{"Jumlah Total Barang", total_product});
    w.Write([]string{"Total Nilai", total_amount_product});
    w.Write([]string{"",});
    w.Write([]string{"SKU", "Nama Barang","Jumlah","Rata Rata Harga Beli","Total"});


    // get available stock
    rows, _ := conn.Db.Raw(sqlquery).Rows()
    defer rows.Close()
    for rows.Next() {
        var m_product_id,name,orderline_qty,average_price,orderline_total_amount string
        rows.Scan(&m_product_id,&name,&orderline_qty,&average_price,&orderline_total_amount)
        var record []string
        record = append(record, m_product_id)
        record = append(record, name)
        record = append(record, orderline_qty)
        record = append(record, average_price)
        record = append(record, orderline_total_amount)
        if err := w.Write(record); err != nil {
            log.Fatalln("error writing record to csv:", err)
        }
    }
    w.Flush()

    if err := w.Error(); err != nil {
        log.Fatal(err)
    }
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Disposition", "attachment; filename=Valueofproduct.csv")
    c.Data(http.StatusOK, "text/csv", b.Bytes())
}

//download report sales order
func SalesOrder(c *gin.Context){

    var date_start string = c.PostForm("date_start")
    var date_end string = c.PostForm("date_end")
    b := &bytes.Buffer{}
    w := csv.NewWriter(b)

    //query get data transaction product
    sqlquery := `
    select *,orderline_total_amount-(average_po_price*orderline_qty) as profit from (
        SELECT
        
            t_orders.t_order_id,
            t_orders.order_date,
            m_products.m_product_id,
            m_products.name,
            t_order_lines.orderline_qty,
            t_order_lines.orderline_price,
            t_order_lines.orderline_total_amount,
            (   
                select 
                    avg(tol.orderline_price) as average_price
                from
                    t_order_lines tol
                    left join t_orders too on (too.t_order_id = tol.t_order_id)
                    where too.order_type = 'P' and tol.m_product_id = t_order_lines.m_product_id
            ) as average_po_price

        from 
            t_order_lines
            left join t_orders on (t_orders.t_order_id = t_order_lines.t_order_id)
            left join m_products on (m_products.m_product_id = t_order_lines.m_product_id)
            where t_orders.order_type = 'S'
            and t_orders.order_date between '`+date_start+`' and '`+date_end+`'
        ) TBX
    `

    //prepare header data
    timenow := time.Now().Format("02-01-2006")
    var total_profit,total_product,omset,total_order string
    
    //subselect query to get summary data
    conn.Db.Raw(
        `select 
            sum(orderline_qty) as total_product,
            sum(orderline_total_amount) as omset, 
            sum(profit) as total_profit
        from (
            `+sqlquery+`
        ) as TB`).Row().Scan(&total_product,&omset,&total_profit)

    //get total order
    conn.Db.Raw(
        `select 
            count(*) as total_order
        from (
            `+sqlquery+`
        ) as TB group by t_order_id`).Row().Scan(&total_order)
    //header csv
    if err := w.Write([]string{"Laporan Nilai Barang"}); err != nil {
            log.Fatalln("error writing record to csv:", err)
    }
    w.Write([]string{"",});
    w.Write([]string{"Tanggal Cetak", timenow});
    w.Write([]string{"Tanggal", date_start+" - "+date_end});
    w.Write([]string{"Total Omset", omset});
    w.Write([]string{"Total Laba Kotor", total_profit});
    w.Write([]string{"Total Penjualan", total_order});
    w.Write([]string{"Total Barang", total_product});
    w.Write([]string{"",});
    

    w.Write([]string{"ID Pesanan", "Waktu","SKU","Nama Barang","Jumlah","Harga Jual","Total","Harga Beli","Laba"});

    // get sales order data
    rows, _ := conn.Db.Raw(sqlquery).Rows()
    defer rows.Close()
    for rows.Next() {
        var t_order_id,m_product_id,name,orderline_qty,orderline_price,orderline_total_amount,average_po_price,profit string
        var order_date time.Time 
        rows.Scan(&t_order_id,&order_date,&m_product_id,&name,&orderline_qty,&orderline_price,&orderline_total_amount,&average_po_price,&profit)
        var record []string
        record = append(record, t_order_id)
        record = append(record, order_date.Format("2006-01-02 15:04:05"))
        record = append(record, m_product_id)
        record = append(record, name)
        record = append(record, orderline_qty)
        record = append(record, orderline_price)
        record = append(record, orderline_total_amount)
        record = append(record, average_po_price)
        record = append(record, profit)
        if err := w.Write(record); err != nil {
            log.Fatalln("error writing record to csv:", err)
        }
    }
    w.Flush()

    if err := w.Error(); err != nil {
        log.Fatal(err)
    }
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Disposition", "attachment; filename=Salesorder.csv")
    c.Data(http.StatusOK, "text/csv", b.Bytes())
}
