package conn
import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)
var (
    // DBCon is the connection handle for the database
	Db *gorm.DB
)

func DbInit() (){
	//open a db connection now is only support sqlite and mysql
	var err error
	dbdriver := viper.GetString("database")
	if dbdriver=="sqlite"{
		Db, err = gorm.Open("sqlite3", "db/inventory.db")
	}else if dbdriver=="mysql"{
		Db, err = gorm.Open("mysql", "root:@/ssinventory?charset=utf8&parseTime=True&loc=Local")
	}
	log.Print("try to connect to ",dbdriver)
	
	if err != nil {
		panic("failed to connect database")
	}
	log.Println("connected to database")
}

