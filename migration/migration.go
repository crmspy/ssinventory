package migration
import (
    "github.com/crmspy/ssinventory/library/conn"
	log "github.com/Sirupsen/logrus"
    "github.com/fatih/structs"
    inventory "github.com/crmspy/ssinventory/modules/inventory/models"
    order "github.com/crmspy/ssinventory/modules/order/models"
)

/*
automatic migrate n seed data from model
*/

func Run(){
    //migrate product
    conn.Db.AutoMigrate(&inventory.Mproduct{})
    conn.Db.AutoMigrate(&order.Torder{},&order.TorderLine{})
    conn.Db.AutoMigrate(&inventory.Minventory{},&inventory.MinventoryLine{},&inventory.Tinout{})
    
    

    model_Minventory := inventory.Minventory{
        M_inventory_id: "General",
        Name: "General Inventory Location",
    }
    conn.Db.Save(&model_Minventory)
    
}

func migrate(table interface{}){
    names := structs.Names(table)
    for _,t := range names {
        log.Println(t)
    }
}
