package main

import (
	"log"

	"github.com/Pumahawk/gjdwp/lib/jdwp"
)

func main() {
	c, err := jdwp.Dial("localhost:5005")
	pack, _ := c.Next()
	log.Printf("pack=%v", pack)
	if err != nil {
		log.Fatalf("%s", err)
	}
}
