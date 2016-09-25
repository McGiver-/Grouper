package handlers

import (
	"fmt"
	"github.com/McGiver-/Grouper/src/myDB"
	"net/http"
)

func GetStats(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("GetStats visited.")
	db := myDB.Db
	stats := db.Stats()
	fmt.Printf("Number of Open Connections are: %d \n", stats.OpenConnections)
}
