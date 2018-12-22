package controller
// create tabel in database if not exist
import (
	"github.com/crmspy/ssinventory/library/conn"
)
func MigrateDb(){
    conn.Db.AutoMigrate(&mProduct{},)
    conn.Db.AutoMigrate(&tOrder{},&tOrderLine{})
    conn.Db.AutoMigrate(&mInventory{},&mInventoryLine{},&tInout{})
}
