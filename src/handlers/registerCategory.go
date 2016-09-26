package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/McGiver-/Grouper/src/myDB"
	"net/http"
)

func RegisterCategory(rw http.ResponseWriter, req *http.Request) {

}

func sendRegisterCategoryResponse(rw *http.ResponseWriter, response CreateCategoryResponse) bool {

	(*rw).Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(*rw).Encode(response); err != nil {
		fmt.Println("Error during json encoding in CreateCategory")
		return false
	} else {
		fmt.Println("Json sent successfully in CreateCategory")
		return true
	}
}
