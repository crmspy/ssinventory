package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
    "github.com/crmspy/ssinventory/library/conn"
    "time"
    "bytes"
    log "github.com/Sirupsen/logrus"
    "encoding/csv"
    
)

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
        m_product.m_product_id,
        m_product.name,
        m_inventory_line.qty_count,
        m_inventory_line.m_inventory_id,
        m_inventory_line.last_update
    from 
        m_inventory_line 
        left join m_product on (m_product.m_product_id = m_inventory_line.m_product_id) 
        order by m_inventory_line.m_product_id asc
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
            t_inout.inout_date,
            m_product.m_product_id,
            m_product.name,
            t_order_line.orderline_qty,
            t_order_line.orderline_received,
            t_order_line.orderline_price,
            t_order_line.orderline_total_amount,
            t_order.t_order_id,
            t_inout.description
        from 
            t_inout
            left join t_order_line on (t_order_line.t_order_line_id = t_inout.t_order_line_id)
            left join m_product on (m_product.m_product_id = t_order_line.m_product_id)
            left join t_order on (t_order.t_order_id = t_order_line.t_order_id)
            where t_inout.inout_type = 'IN'
            group by t_inout.t_order_line_id
            order by t_inout.inout_date asc
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
            t_inout.inout_date,
            m_product.m_product_id,
            m_product.name,
            t_order_line.orderline_qty,
            t_order_line.orderline_received,
            t_order_line.orderline_price,
            t_order_line.orderline_total_amount,
            t_order.t_order_id,
            t_inout.description
        from 
            t_inout
            left join t_order_line on (t_order_line.t_order_line_id = t_inout.t_order_line_id)
            left join m_product on (m_product.m_product_id = t_order_line.m_product_id)
            left join t_order on (t_order.t_order_id = t_order_line.t_order_id)
            where t_inout.inout_type = 'OUT'
            group by t_inout.t_order_line_id
            order by t_inout.inout_date asc
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
        m_product.m_product_id,
        m_product.name,
        sum(orderline_qty) as orderline_qty,
        avg(orderline_price) as average_price,
        sum(orderline_total_amount) as orderline_total_amount
    from
        t_order_line
        left join m_product on (m_product.m_product_id = t_order_line.m_product_id)
        left join t_order on (t_order.t_order_id = t_order_line.t_order_id)
        where t_order.order_type = 'P'
        GROUP by t_order_line.m_product_id
        order by t_order_line.m_product_id asc
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
        
            t_order.t_order_id,
            t_order.order_date,
            m_product.m_product_id,
            m_product.name,
            t_order_line.orderline_qty,
            t_order_line.orderline_price,
            t_order_line.orderline_total_amount,
            (   
                select 
                    avg(tol.orderline_price) as average_price
                from
                    t_order_line tol
                    left join t_order too on (too.t_order_id = tol.t_order_id)
                    where too.order_type = 'P' and tol.m_product_id = t_order_line.m_product_id
            ) as average_po_price

        from 
            t_order_line
            left join t_order on (t_order.t_order_id = t_order_line.t_order_id)
            left join m_product on (m_product.m_product_id = t_order_line.m_product_id)
            where t_order.order_type = 'S'
            and t_order.order_date between '`+date_start+`' and '`+date_end+`'
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
