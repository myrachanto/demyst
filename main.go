package main

import (
	"log"

	"github.com/myrachanto/estate/src/routes"
)

func init() {
	log.SetPrefix("Demyst server ...... ")
}
func main() {
	log.Println("Server started")
	routes.ApiLoader()
}
