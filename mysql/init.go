package mysql

import (
	"fmt"
	"log"
	"go-restapi/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func init() {
	config := config.GetConfig()
	
	path := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", config.DatabaseUser, config.DatabasePassword, config.DatabaseHost, config.DatabasePort, config.DatabaseName)
	
	var err error
	if Db, err = sql.Open("mysql", path); err != nil {
		log.Fatalf("Failed to open mysql: %v", err)
	}

	log.Println("Successfully connected to mysql")
}