package main

import (
	"github.com/davecgh/go-spew/spew"
	"gobserver/internal/config"
	"log"
)

func main() {
	appConfig, err := config.NewConfig("./", "config")
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(appConfig)
}
