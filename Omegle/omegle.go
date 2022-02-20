package omegle

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Omegle struct {
	clientID    string
	client      *http.Client
	lock        *sync.Mutex //resource sharing try not cause a data race
	wg          *sync.WaitGroup
	serverurl   string   // server url
	start       string   // init url link
	wpm         int      //words per min
	firstevents int      //random shit again
	server      int      // concat with a formatted string later
	loop        bool     //loop again?
	topics      []string //slice of topics to include
	events      Events   //struct of events
	lang        string   //language!
	likesvalid  bool
	connected   bool
	quit        chan bool
}

type Options struct {
	Wpm     int
	Server  int
	Loop    bool
	Topics  []string
	Timeout int
}

func (o *Omegle) ChkLikesvalid() bool {
	return o.likesvalid
}

func (o *Omegle) ChkConnected() bool {
	return o.connected
}

func (o *Omegle) ChangeLang(newlang string) { //fix later possible data race may occur here
	o.lock.Lock()
	o.lang = newlang
	o.lock.Unlock()
}

func (o *Omegle) ChangeWpm(newwpm int) {
	o.lock.Lock()
	o.wpm = newwpm
	o.lock.Unlock()
}

func (o *Omegle) ChangeLoop(newloop bool) {
	o.lock.Lock()
	o.loop = newloop
	o.lock.Unlock()
}

func (o *Omegle) randID(length int) (Id string) {
	for i := 0; i < length; i++ {
		Id += string(genid[rand.Intn(len(genid))])
	}
	return Id
}

func (o *Omegle) genLink() {
	o.start = fmt.Sprintf(START_URL, o.serverurl, o.firstevents, o.randID(5), o.lang)
	if len(o.topics) > 0 {
		topicsdata, err := json.Marshal(o.topics)
		if err != nil {
			fmt.Printf("An error occured while marshalling json topics: %v\n", err)
			return
		}
		o.start += "&topics=" + url.QueryEscape(string(topicsdata))
		fmt.Println(o.start)
	}
}

func (o *Omegle) Run() {
	o.genLink()
	o.Connect()
	o.wg.Add(1)
	go o.listener()
	o.wg.Wait()
}

func (o *Omegle) Exit() {
	if o.connected {
		err := o.Disconnect()
		if err != nil {
			fmt.Printf("An error occurred while exiting in disconnect: %v\n", err)
		}
	}
	o.quit <- true //send a signal to listener to quit

}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewOmegle(options Options, events Events) *Omegle {
	if options.Wpm == 0 {
		options.Wpm = 42
	}
	if options.Server == 0 {
		options.Server = rand.Intn(47) + 1
	}
	if options.Timeout == 0 {
		options.Timeout = 15
	}
	checkEvents(&events)
	return &Omegle{
		client: &http.Client{
			Timeout: time.Duration(options.Timeout) * time.Second,
		},
		lock:      &sync.Mutex{},
		wg:        &sync.WaitGroup{},
		serverurl: fmt.Sprintf(SERVER_URL, options.Server),
		wpm:       options.Wpm,
		server:    options.Server,
		loop:      options.Loop,
		topics:    options.Topics,
		events:    events,
	}

}
