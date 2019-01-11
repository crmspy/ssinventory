package conn
import (
	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    _ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/Sirupsen/logrus"
    "github.com/spf13/viper"
)
var (
    // DBCon is the connection handle for the database
	Db *gorm.DB
)


func DbInit() (){
    /*
    open a db connection now is only support sqlite and mysql
    this configuration read from config/app.json file
    */
	var err error
    dbdriver := viper.GetString("database")
	if dbdriver=="sqlite"{
		Db, err = gorm.Open("sqlite3", viper.GetString("sqlite.file"))
	}else if dbdriver=="mysql"{
        var host string = viper.GetString("mysql.host")
        var port string = viper.GetString("mysql.port")
        var username string = viper.GetString("mysql.username")
        var password string = viper.GetString("mysql.password")
        var database string = viper.GetString("mysql.database")

		Db, err = gorm.Open("mysql", username+":"+password+"@tcp("+host+":"+port+")/"+database+"?charset=utf8&parseTime=True&loc=Local")
	}else if dbdriver=="postgresql"{
        /*
            Work In Progress 
            table success imported is m_inventories and m_product other that failed to migrate
        */
        var host string = viper.GetString("postgresql.host")
        var port string = viper.GetString("postgresql.port")
        var username string = viper.GetString("postgresql.username")
        var password string = viper.GetString("postgresql.password")
        var database string = viper.GetString("postgresql.database")
        var sslmode string = viper.GetString("postgresql.sslmode")
        Db, err = gorm.Open("postgres", "host="+host+" port="+port+" user="+username+" dbname="+database+" password="+password+" sslmode="+sslmode+"")
        
    }
	log.Print("try to connect to ",dbdriver)
	
	if err != nil {
        log.Println(err)
		panic("failed to connect database")
	}
	log.Println("connected to database")
}

