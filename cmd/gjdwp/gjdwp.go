package main

import (
	"log"

	"github.com/Pumahawk/gjdwp/lib/jdwp"
)

func main() {
	_, err := jdwp.Dial("localhost:5005")
	if err != nil {
		log.Fatalf("%s", err)
	}
}
