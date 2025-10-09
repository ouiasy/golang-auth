package main

import (
	"log"
)

func main() {
	if err := serve(); err != nil {
		log.Fatalln("error running server: ", err.Error())
	}
}
