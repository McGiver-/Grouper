package myDB

import (
	"fmt"
	"database/sql"
)

type MyDB struct{
	db *sql.DB
}

func Inititialize(driverName, dataSourceName string) (mydb MyDB){
	var db *sql.DB
	db, err = sql.Open(driverName, dataSourceName)

	if err != nil {
		fmt.Printf("Connected to database at $1",dataSourceName)
	}else{
		mydb = MyDB{db}
	}
	return mydb, err
}

