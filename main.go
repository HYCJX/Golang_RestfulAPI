package main

import (
	"fmt"
	"github.com/HYCJX/Golang_RestfulAPI/api"
	"log"
)
func main() {
	crew, err := api.ReadJson()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*crew)
}
