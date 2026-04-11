package main

import (
	"fmt"
	"log"

	"github.com/Pumahawk/gjdwp/lib/jdwp"
)

func main() {
	c, err := jdwp.Dial("localhost:5005")
	if err != nil {
		log.Fatalf("%s", err)
	}
	pack, _ := c.Next()
	switch pack := pack.(type) {
	case *jdwp.ReplyPack:
		log.Printf("reply pack=%v\n", pack)
	case *jdwp.CommandPack:
		switch event := pack.Data.(type) {
		case *jdwp.Event:
			for i, eventData := range event.EventData {
				i := i + 1
				switch eventData.(type) {
				case *jdwp.VMStartEventData:
					log.Printf("command pack: event: (%d/%d) vmstartevent SuspendPolicy=%q\n", i, event.Events, event.SuspendPolicy)
				default:
					log.Printf("command pack: event: (%d/%d) unknown data type %T\n", i, event.Events, event.EventData)
				}
			}
		default:
			log.Printf("command pack: unknown data type %T\n", pack.Data)
		}
	default:
		panic(fmt.Errorf("main unknown pack type %T", pack))
	}

	log.Printf("Write ResumeCommandData")
	err = c.Write(jdwp.NewResumeCommand(102))
	if err != nil {
		log.Printf("Error on write pack: %s", err)
	} else {
		log.Printf("Success write pack")
	}
	ch := make(chan int)
	<-ch
}
