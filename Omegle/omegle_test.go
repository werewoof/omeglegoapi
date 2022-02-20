package omegle

import (
	"fmt"
	"testing"
)

func TestClientConnectAndDisconnect(t *testing.T) {
	options := Options{
		Wpm:     42,
		Timeout: 15,
	}
	events := Events{
		On_ready: func(o *Omegle) {
			fmt.Println("on_ready passed")
		},
		On_message: func(o *Omegle, message string) {
			fmt.Println("Stranger said:", message)
		},
	}
	client := NewOmegle(options, events)
	client.Run()
}
