package autoload

/*
load all configuration
static parameter n database connection
*/

import (
	"os"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/crmspy/ssinventory/library/conn"
    "github.com/crmspy/ssinventory/migration"
)

func Run() {
	logger()
	//banner()
	readConfig()
    conn.DbInit()
    migration.Run()
}

// load application configuration
func readConfig() {
	log.Println("Initialing File Configuration")

	//application configuration
	viper.SetConfigFile("./config/app.json")
	// Searches for config file in given paths and read it
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Fatal Error reading config file, %s", err)
	}

}

//activate loging
func logger() {
	var filename string = "./log/runtime.log"
	// Create the log file if doesn't exist. And append to it if it already exists.
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	if err != nil {
		// Cannot open log file. Logging to stderr
        log.Fatalf("Cannot open Log file", err)
	} else {
		log.SetOutput(f)
	}
}

//show banner data
func banner() {
	var header string = `
	.__   .__           ___.    
________ __ __ |  |  |  |  _____   \_ |__  
\___   /|  |  \|  |  |  |  \__  \   | __ \ 
/    / |  |  /|  |__|  |__ / __ \_ | \_\ \
/_____ \|____/ |____/|____/(____  / |___  /
\/                        \/      \/ 
============================================
`
	log.Println(header)
	log.Println("Starting application")
}
