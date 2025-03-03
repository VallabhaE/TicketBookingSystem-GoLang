package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var (
	database *sql.DB
	
)

// Init is Mandatory before moving to other functions
func InitDB(){
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s","root","root","MovieTicketBooking")
	var err error
	database, err = sql.Open("mysql", dsn)
	if err!=nil{
		panic(err)
	}
}

func Ping(){
	err := database.Ping()
	if err!=nil{
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Pong")
}


// To Get Db Object and Use Available Inbuilt functions to the database
func GetDbObject() *sql.DB{
	return database
}