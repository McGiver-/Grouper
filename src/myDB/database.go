package myDB

import (
	"fmt"
	"database/sql"
	"log"
)

var Db *sql.DB

func Init(driverName, dataSourceName string){
	Db, err = sql.Open(driverName, dataSourceName)

	if err != nil {
		fmt.Printf("***Connection to database failed at $1 \n",dataSourceName)
		log.Fatal("Could Not Connect To database")
	}else{
		fmt.Println("Connected to database")
	}
}
