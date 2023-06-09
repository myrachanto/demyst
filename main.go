package main

import (
	"log"

	"github.com/myrachanto/sports/src/routes"
)

func init() {
	log.SetPrefix("Demyst server ...... ")
}
func main() {
	log.Println("Server started")
	routes.ApiLoader()
}
