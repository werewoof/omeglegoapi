package omegle

import "fmt"

type Events struct {
	On_ready          func(o *Omegle) //pass object for now might change how this works later
	On_wait           func(o *Omegle) //expliclity passing object itself not sure if legal
	On_message        func(o *Omegle, message string)
	On_typing         func(o *Omegle)
	On_stopped_typing func(o *Omegle)
	On_disconnect     func(o *Omegle)

	On_idle func(o *Omegle) //when the omegle client is disconnected

	//Not really important stuff....

	RecaptchaRequired func(o *Omegle)
	RecaptchaRejected func(o *Omegle)
	Common_likes      func(o *Omegle, likes string)
	ConnectionDied    func(o *Omegle)
	IdentDigests      func(o *Omegle, message string)
	StatusInfo        func(o *Omegle, statusinfo Status)
	ServerMessage     func(o *Omegle, message string)
	AntinudeBanned    func(o *Omegle)
}

func checkEvents(events *Events) { //replaces all nil functions with empty ones
	if events.On_ready == nil {
		events.On_ready = func(o *Omegle) {}
	}
	if events.On_wait == nil {
		events.On_wait = func(o *Omegle) {}
	}
	if events.On_message == nil {
		events.On_message = func(o *Omegle, message string) {}
	}
	if events.On_typing == nil {
		events.On_typing = func(o *Omegle) {}
	}
	if events.On_stopped_typing == nil {
		events.On_stopped_typing = func(o *Omegle) {}
	}
	if events.On_disconnect == nil {
		events.On_disconnect = func(o *Omegle) {}
	}
	if events.RecaptchaRejected == nil {
		events.RecaptchaRejected = func(o *Omegle) {}
	}
	if events.Common_likes == nil {
		events.Common_likes = func(o *Omegle, likes string) {}
	}
	if events.ConnectionDied == nil {
		events.ConnectionDied = func(o *Omegle) {}
	}
	if events.IdentDigests == nil {
		events.IdentDigests = func(o *Omegle, message string) {}
	}
	if events.StatusInfo == nil {
		events.StatusInfo = func(o *Omegle, statusinfo Status) {}
	}
	if events.ServerMessage == nil {
		events.ServerMessage = func(o *Omegle, message string) {}
	}
	if events.On_idle == nil {
		events.On_idle = func(o *Omegle) {}
	}
}

//supposed to look like ["gotMessage", "hello!"]
func (o *Omegle) event_selector(events []interface{}) error {
	switch event := events[0].(string); event {
	case "waiting":
		o.wg.Add(1)
		go func() { //try find another way next time
			o.events.On_wait(o)
			o.wg.Done()
		}()
	case "typing":
		o.wg.Add(1)
		go func() {
			o.events.On_typing(o)
			o.wg.Done()
		}()
	case "connected":
		o.wg.Add(1)
		go func() {
			o.events.On_ready(o)
			o.wg.Done()
		}()
	case "gotMessage":
		o.wg.Add(1)
		go func() {
			o.events.On_message(o, events[1].(string))
			o.wg.Done()
		}()
	case "commonLikes":
		o.likesvalid = true
		o.wg.Add(1)
		go func() {
			o.events.Common_likes(o, events[1].(string))
			o.wg.Done()
		}()
	case "stoppedTyping":
		o.wg.Add(1)
		go func() {
			o.events.On_stopped_typing(o)
			o.wg.Done()
		}()
	case "strangerDisconnected":
		o.wg.Add(1)
		go func() {
			o.events.On_disconnect(o)
			if o.loop {
				o.Connect()
			}
			o.wg.Done()
		}()
	case "recaptchaRequired":
		o.wg.Add(1)
		go func() {
			o.events.RecaptchaRequired(o)
			o.wg.Done()
		}()
	case "recaptchaRejected":
		o.wg.Add(1)
		go func() {
			o.events.RecaptchaRejected(o)
			o.wg.Done()
		}()
	case "connectionDied":
		o.wg.Add(1)
		go func() {
			o.events.ConnectionDied(o)
			o.wg.Done()
		}()
	case "identDigests":
		o.wg.Add(1)
		go func() {
			o.events.IdentDigests(o, events[1].(string))
			o.wg.Done()
		}()
	case "statusInfo":
		o.wg.Add(1)
		go func() {
			o.events.StatusInfo(o, events[1].(Status))
			o.wg.Done()
		}()
	case "serverMessage":
		o.wg.Add(1)
		go func() {
			o.events.ServerMessage(o, events[1].(string))
			o.wg.Done()
		}()
	case "antinudeBanned":
		o.wg.Add(1)
		go func() {
			o.events.AntinudeBanned(o)
			o.wg.Done()
		}()
	default:
		return fmt.Errorf("unhandled event type %s", event)
	}
	return nil
}
