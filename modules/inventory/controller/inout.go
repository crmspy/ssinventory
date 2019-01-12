package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
    "fmt"
    "strconv"
    "github.com/crmspy/ssinventory/library/conn"
    "time"
    "errors"
    log "github.com/Sirupsen/logrus"
)


type    InoutParam struct {
    Inout_type		        string
    M_inventory_id			string
    M_product_id			string
    T_order_line_id         int
    Inout_qty               int
    Description             string
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
    param := InoutParam {
        M_inventory_id: c.PostForm("m_inventory_id"),
        Inout_qty: qty,
        T_order_line_id: t_order_line_id,
        Description: c.PostForm("description"),
    }
    if e := DoInout(param); e != nil{
        c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprint(e)})
    }else{
        c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "inventory '"+param.M_inventory_id+"' was updated successfully"})
    }
}

/*
the main function of the inventory application
Automatically update stock data in the warehouse and record the history of product IN and OUT
Automatically detect it's sales order / purchase order

IN  = Put the product into the warehouse
OUT = Eject product from the warehouse
P   = Purchase Order
S   = Sales Order
*/
func DoInout(p InoutParam)(error){
    if p.M_inventory_id=="" || p.T_order_line_id ==0 || p.Inout_qty==0 {
        return errors.New("Please fill mandatory parameter")
    }

    var M_inventory_id string
    conn.Db.Table("m_inventory").Where("m_inventory_id = ?", p.M_inventory_id).Select("m_inventory_id").Row().Scan(&M_inventory_id)

    if M_inventory_id=="" {
        return errors.New("Cannot find inventory location in database")
    }

    var order_type,m_product_id string
    var orderline_qty,orderline_outstanding,orderline_received int

    conn.Db.Raw("SELECT order_type,orderline_qty,orderline_outstanding,orderline_received,m_product_id FROM t_order_line join t_order using (t_order_id) WHERE t_order_line_id = ?", p.T_order_line_id).Row().Scan(&order_type,&orderline_qty,&orderline_outstanding,&orderline_received,&m_product_id)
    if m_product_id==""{
        return errors.New("Transaction not found")
    }else if orderline_outstanding==0 {
        return errors.New("Sorry, all product "+m_product_id +" from this transaction has been received")
    }else if orderline_outstanding < p.Inout_qty {
        return errors.New("maximum received product "+m_product_id +" from this transaction is "+string(orderline_outstanding)+" qty")
    }

    product_received    := orderline_received + p.Inout_qty
    product_outstanding := orderline_outstanding - p.Inout_qty
    timenow := time.Now().Format("2006-01-02 15:04:05")

    log.Println("Restock Product ",p)
    switch order_type {
    case "P":
        //purchase order
        //insert into tinout
        conn.Db.Exec("INSERT INTO t_inout (inout_type, m_inventory_id, m_product_id, t_order_line_id, inout_qty,inout_date,description) VALUES(?, ?, ?, ?, ?, ?, ?)", "IN",p.M_inventory_id,m_product_id,p.T_order_line_id,p.Inout_qty,timenow,p.Description)
        
        //update data order
        conn.Db.Exec("UPDATE t_order_line SET orderline_outstanding=?, orderline_received=? where t_order_line_id = ?", product_outstanding, product_received,p.T_order_line_id)
        
        //insert into tinout for history movement product
        var qty_count,m_inventory_line_id int
        conn.Db.Raw("SELECT m_inventory_line_id,qty_count FROM m_inventory_line WHERE m_product_id = ? and m_inventory_id = ?", m_product_id,p.M_inventory_id).Row().Scan(&m_inventory_line_id,&qty_count)
        
        if (m_inventory_line_id > 0){
            
            //update stock
            new_stock := qty_count + p.Inout_qty
            conn.Db.Exec("UPDATE m_inventory_line SET qty_count=?, last_update=? where m_inventory_line_id = ?", new_stock, timenow,m_inventory_line_id)

        }else{
            
            //insert new stock
            new_stock := p.Inout_qty
            conn.Db.Exec("INSERT INTO m_inventory_line (m_inventory_id, m_product_id, qty_count,last_update) VALUES(?, ?, ?, ?)",p.M_inventory_id,m_product_id,new_stock,timenow)
        
        }
    case "S":
        //Sales Order
        
        //check stock is available or not
        var qty_count,m_inventory_line_id int
        conn.Db.Raw("SELECT m_inventory_line_id,qty_count FROM m_inventory_line WHERE m_product_id = ? and m_inventory_id = ?", m_product_id,p.M_inventory_id).Row().Scan(&m_inventory_line_id,&qty_count)
        if qty_count < p.Inout_qty {
            return errors.New("Sorry "+m_product_id+" out of stock, only available "+string(qty_count)+" qty")
        }
        
        //insert into tinout for history movement product
        conn.Db.Exec("INSERT INTO t_inout (inout_type, m_inventory_id, m_product_id, t_order_line_id, inout_qty,inout_date,description) VALUES(?, ?, ?, ?, ?, ?, ?)", "OUT",p.M_inventory_id,m_product_id,p.T_order_line_id,p.Inout_qty,timenow,p.Description)
        
        //update data order
        conn.Db.Exec("UPDATE t_order_line SET orderline_outstanding=?, orderline_received=? where t_order_line_id = ?", product_outstanding, product_received,p.T_order_line_id)
        if (m_inventory_line_id > 0){
            
            //update stock
            new_stock := qty_count - p.Inout_qty
            conn.Db.Exec("UPDATE m_inventory_line SET qty_count=?, last_update=? where m_inventory_line_id = ?", new_stock, timenow,m_inventory_line_id)

        }else{
            
            //insert new stock
            new_stock := p.Inout_qty
            conn.Db.Exec("INSERT INTO m_inventory_line (m_inventory_id, m_product_id, qty_count,last_update) VALUES(?, ?, ?, ?)",p.M_inventory_id,m_product_id,new_stock,timenow)
        
        }
    default:
        return errors.New("Inout Type not found")
    }

    return nil
}
