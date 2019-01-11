package models

import (
    "time"
)
//change table name
func (Torder) TableName() string {
    return "t_order"
}

//change table name
func (TorderLine) TableName() string {
    return "t_order_line"
}
type (
	Torder struct {
		//gorm.Model
		T_order_id		string		`gorm:"type:varchar(64);PRIMARY_KEY"`
		Order_type		string		`gorm:"type:varchar(1);"`
        Order_amount	float64		`gorm:"type:float(16,2);"`
        Description     string
        Order_status	string		`gorm:"type:varchar(1);"`
        Order_date      time.Time   `gorm:"type:datetime"`
    }

    TorderLine struct {
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
	TransformedTorder struct {
		T_order_id		string		`json:"transaction_id"`
		Order_type		string		`json:"order_type"`
        Order_amount	float64		`json:"amount"`
        Description     string      `json:"description"`
        Order_status    string      `json:"order_status"`
        Order_date      time.Time   `json:"order_date"`
    }
    
    TransformedTorderLine struct {
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
