package omegle

import (
	"encoding/json"
	"fmt"
	"time"
)

func (o *Omegle) EmuTyping(message string) error {
	seconds := time.Duration((60 / o.wpm) * (len(message) / 5))
	err := o.Typing()
	if err != nil {
		return err
	}
	time.Sleep(seconds * time.Second)
	err = o.StopTyping()
	if err != nil {
		return err
	}
	err = o.Send(message)
	if err != nil {
		return err
	}
	return nil
}

func (o *Omegle) Typing() error {
	if !o.connected {
		return fmt.Errorf("not connected to any chat")
	}
	data := IdData{
		o.clientID,
	}
	response, err := o.PostRequest(fmt.Sprintf(TYPING_URL, o.serverurl), data)
	if response != "win" {
		return fmt.Errorf("post request to %s is not successful", fmt.Sprintf(TYPING_URL, o.serverurl))
	}
	if err != nil {
		return err
	}
	return nil
}

func (o *Omegle) StopTyping() error {
	if !o.connected {
		return fmt.Errorf("not connected to any chat")
	}
	data := IdData{
		o.clientID,
	}
	response, err := o.PostRequest(fmt.Sprintf(STOP_TYPING_URL, o.serverurl), data)
	if response != "win" {
		return fmt.Errorf("post request to %s is not successful", fmt.Sprintf(STOP_TYPING_URL, o.serverurl))
	}
	if err != nil {
		return err
	}
	return nil
}

func (o *Omegle) Send(message string) error {
	if !o.connected {
		return fmt.Errorf("not connected to any chat")
	}
	data := IdMsgdata{
		o.clientID,
		message,
	}
	response, err := o.PostRequest(fmt.Sprintf(SEND_URL, o.serverurl), data)
	if response != "win" {
		return fmt.Errorf("post request to %s is not successful", fmt.Sprintf(SEND_URL, o.serverurl))
	}
	if err != nil {
		return err
	}
	return nil
}

func (o *Omegle) Connect() error {
	response, err := o.PostRequest(o.start, Nothing{})
	if err != nil {
		return err
	}
	var startResp StartResp
	err = json.Unmarshal([]byte(response), &startResp)
	if err != nil {
		return err
	}
	o.lock.Lock()
	o.clientID = startResp.ClientID
	o.connected = true
	o.likesvalid = false
	o.lock.Unlock()
	for _, event := range startResp.Events {
		err := o.event_selector(event)
		if err != nil {
			fmt.Printf("An Error occurred while reading events from connection: %v\n", err)
		}
	}
	return nil
}

func (o *Omegle) Disconnect() error {
	data := IdData{
		o.clientID,
	}
	response, err := o.PostRequest(fmt.Sprintf(DISCONNECT_URL, o.serverurl), data)
	if err != nil {
		return err
	}
	if response != "win" {
		return fmt.Errorf("post request to %s is not successful", fmt.Sprintf(DISCONNECT_URL, o.serverurl))
	}
	o.lock.Lock()
	o.connected = false
	o.likesvalid = false
	o.clientID = ""
	o.lock.Unlock()
	return nil
}

func (o *Omegle) StopLikes() error {
	data := IdData{
		o.clientID,
	}
	response, err := o.PostRequest(fmt.Sprintf(STOPLIKES_URL, o.serverurl), data)
	if err != nil {
		return err
	}
	if response != "win" {
		return fmt.Errorf("post request to %s is not successful", fmt.Sprintf(STOPLIKES_URL, o.serverurl))
	}
	return nil
}

func (o *Omegle) GetStatus() (Status, error) {
	data := IdData{
		o.clientID,
	}
	response, err := o.GetRequest(fmt.Sprintf(STATUS_URL, o.serverurl, o.randID(5)), data)
	if err != nil {
		return Status{}, err
	}
	var status Status
	err = json.Unmarshal([]byte(response), &status)
	if err != nil {
		return Status{}, err
	}
	return status, nil
}

func (o *Omegle) listener() { //LONG POLLING returning null
	fmt.Println("passed")
	for {
		select {
		case <-o.quit: // quit listener once client decides to exit
			o.wg.Done()
			return
		default:
			if !o.connected {
				o.events.On_idle(o) //not async for this function otherwise its gonna spam this function
				continue
			}
			fmt.Println(o.clientID)
			data := IdData{
				o.clientID,
			}
			response, err := o.PostRequest(fmt.Sprintf(EVENT_URL, o.serverurl), data)
			if err != nil {
				fmt.Printf("An Error occurred in the listener: %v\n", err)
			}
			if response == "null" {
				fmt.Println("An error occured in the listener: events returned null")
			}
			fmt.Print(response)
			var events EventsResp
			err = json.Unmarshal([]byte(response), &events)
			fmt.Println(events)
			if err != nil {
				fmt.Printf("An Error occured in the listener: %v\n", err)
			}
			for _, event := range events {
				err := o.event_selector(event) //since I converted all values into interfaces i need to expliclty cast it into a string
				if err != nil {
					fmt.Printf("An Error occurred while reading events from listener: %v\n", err)
				}
			}
		}
	}
}
