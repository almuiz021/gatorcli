package main

import (
	"fmt"
	"log"

	"github.com/almuiz021/gatorcli/internal/config"
)

func main() {
	dbConfig, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", dbConfig)

	err = dbConfig.SetUser("abdulmuiz")
	if err != nil {
		log.Fatalf("couldn't set current user: %v", err)
	}

	dbConfig, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config again: %+v\n", dbConfig)

}
