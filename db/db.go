package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type SQLDB interface {
	SetupDatabaseConnection() *gorm.DB
}

type Sql struct {
}

func NewSQL() SQLDB {
	return &Sql{}
}

func (s *Sql) SetupDatabaseConnection() (db *gorm.DB) {
	err := godotenv.Load() //load env file
	if err != nil {
		log.Fatal(err.Error())
		return db
	}

	//get database config value from env file
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUsername, dbPassword, dbHost, dbPort, dbName)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}) //open connection to mysql drive
	if err != nil {
		log.Fatal(err.Error())

		return db
	}

	return db
}
