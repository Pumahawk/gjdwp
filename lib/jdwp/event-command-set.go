package jdwp

import (
	"bytes"
	"fmt"
)

type SuspendPolicy uint8

func (s SuspendPolicy) String() string {
	switch s {
	case NoneSP:
		return "NoneSP"
	case EventThreadSP:
		return "EventThreadSP"
	case AllSP:
		return "AllSP"
	default:
		return fmt.Sprint(uint8(s))
	}
}

const (
	NoneSP SuspendPolicy = iota
	EventThreadSP
	AllSP
)

type EventData any

type Event struct {
	SuspendPolicy SuspendPolicy
	Events        uint32
	EventData     []EventData
}

type VMStartEventData struct {
	RequestId uint32
	ThreadId  uint64
}

func readEventComposite(bf []byte) (*Event, error) {
	r := bytes.NewReader(bf)
	var suspendPolicy SuspendPolicy
	if err := binRead(r, &suspendPolicy); err != nil {
		return nil, fmt.Errorf("jdwp read event composite suspendPolicy: %w", err)
	}
	var events uint32
	if err := binRead(r, &events); err != nil {
		return nil, fmt.Errorf("jdwp read event composite events: %w", err)
	}
	eventData := make([]EventData, 0, events)
	var eventKind uint8
	for range events {
		if err := binRead(r, &eventKind); err != nil {
			return nil, fmt.Errorf("jdwp read event composite eventKind: %w", err)
		}
		switch eventKind {
		case 90:
			var vmStartRequestData VMStartEventData
			if err := binRead(r, &vmStartRequestData); err != nil {
				return nil, fmt.Errorf("jdwp read event composite events vmStartRequestData: %w", err)
			}
			eventData = append(eventData, &vmStartRequestData)
		default:
			return nil, fmt.Errorf("jdwp read event composite not supported eventKind: missing implementation")
		}
	}
	return &Event{suspendPolicy, events, eventData}, nil
}
