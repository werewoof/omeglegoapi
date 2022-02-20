package main

import (
	"fmt"

	"github.com/asianchinaboi/omegleclient/omegle"
)

func main() {
	options := omegle.Options{
		Wpm:     42,
		Timeout: 15,
		Topics:  []string{"dababybot"},
	}
	events := omegle.Events{
		On_ready: func(o *omegle.Omegle) {
			fmt.Println("on_ready passed")
		},
		On_message: func(o *omegle.Omegle, message string) {
			fmt.Println("Stranger said:", message)
			o.EmuTyping("Hello!")
		},
		Common_likes: func(o *omegle.Omegle, likes string) {
			fmt.Println("You and this stranger have common likes", likes)
		},
	}
	client := omegle.NewOmegle(options, events)
	client.Run()
}
