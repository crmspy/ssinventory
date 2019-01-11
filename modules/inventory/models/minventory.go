package models

import (
    "time"
)

//change table name
func (Minventory) TableName() string {
    return "m_inventory"
}

func (MinventoryLine) TableName() string {
    return "m_inventory_line"
}

func (Tinout) TableName() string {
    return "t_inout"
}

type (
    //inventory location
	Minventory struct {
		M_inventory_id	        string		`gorm:"type:varchar(64);PRIMARY_KEY"`
		Name			        string		`gorm:"type:varchar(255);"`
	}

    //information all product at inventory
	MinventoryLine struct {
		M_inventory_line_id		int			`gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
		M_inventory_id			string		`gorm:"type:varchar(64);"`
		M_product_id			string		`gorm:"type:varchar(64);"`
        Qty_count               int         `gorm:"type:int"`
        Last_update             time.Time   `gorm:"type:datetime"`
	}

    //Inout product at inventory
    Tinout struct {
        T_inout_id		        int			`gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
        Inout_type		        string		`gorm:"type:varchar(3);"`
		M_inventory_id			string		`gorm:"type:varchar(64);"`
        M_product_id			string		`gorm:"type:varchar(64);"`
        T_order_line_id         int         `gorm:"type:int"`
        Inout_qty               int         `gorm:"type:int"`
        Inout_date              time.Time   `gorm:"type:datetime"`
        Description             string      `gorm:"type:varchar(255)"`
    }

    // transformedInventory represents a formatted inventory location
	TransformedMinventory struct {
		M_inventory_id	string		`json:"m_inventory_id"`
		Name			string		`json:"name"`
    }
    
)
