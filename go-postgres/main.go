// psql -d template1 --- to open postgres in terminal
package main

import (
	"fmt"
	"log"
	"net/http"
	"go-postgres/router"
)

func main(){
	r := router.Router()
	fmt.Println("Starting server on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}