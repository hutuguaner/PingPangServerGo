package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"time"
	
	
)





const (
	USERNAME = "root"
	PASSWORD = "root"
	NETWORK  = "tcp"
	SERVER   = "localhost"
	PORT     = 3306
	DATABASE = "pingpang"
)

var hasInitDB bool = false
var DB *sql.DB = nil
var openDBError error = nil
func initDB() {
	hasInitDB = true
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	DB, openDBError = sql.Open("mysql", dsn)
	if openDBError != nil {
		fmt.Printf("open mysql failed,err:%v\n", openDBError)
		if DB != nil {
			DB.Close()
			DB = nil
		}
		return
	}
	DB.SetConnMaxLifetime(100 * time.Second)
	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(16)
}


func main() {

	http.HandleFunc("/getAllUser", getAllUser)
	http.HandleFunc("/login",login)
	http.HandleFunc("/regist",regist)
	http.ListenAndServe(":9090", nil)

}
