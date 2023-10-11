package config

import (
	"os"
)

type Config struct {
	Port              string    
	DatabaseName      string 
	DatabaseUser      string 
	DatabasePassword  string 
	DatabasePort      string 
	DatabaseHost      string
}

//return config struct
func GetConfig() *Config {
	return &Config{
		Port              :os.Getenv("PORT_NUMBER"),  
		DatabaseName      :os.Getenv("MYSQL_DATABASE"),  
		DatabaseUser      :os.Getenv("MYSQL_USER"),  
		DatabasePassword  :os.Getenv("MYSQL_PASSWORD"),   
		DatabasePort      :os.Getenv("MYSQL_PORT"),   
		DatabaseHost      :os.Getenv("MYSQL_DB_HOST"),
	}
}
