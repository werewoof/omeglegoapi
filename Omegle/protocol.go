package omegle

const (
	genid = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
)

const ( //could use %v but i want to make this strict
	SERVER_URL      = "front%d.omegle.com"
	STATUS_URL      = "http://%s/status?randid=%s" //not used
	START_URL       = "http://%s/start?caps=recaptcha2,t&firstevents=%d&randid=%s&lang=%s"
	EVENT_URL       = "http://%s/events"
	TYPING_URL      = "http://%s/typing"
	SEND_URL        = "http://%s/send"
	STOP_TYPING_URL = "http://%s/stoptyping"
	DISCONNECT_URL  = "http://%s/disconnect"
	STOPLIKES_URL   = "https://%s/stoplookingforcommonlikes"
)

type Status struct {
	Count           int      `json:"count"`
	Force_unmon     bool     `json:"force_unmon"`
	Antinudeservers []string `json:"antinudeservers"`
	Antinudepercent float32  `json:"antinudepercent"`
	SpyQueueTime    float32  `json:"spyQueueTime"`
	TimeStamp       float32  `json:"timestamp"`
	Servers         []string `json:"servers"`
}

type StartResp struct {
	Events     EventsResp `json:"events"`
	ClientID   string     `json:"clientID"`
	StatusInfo Status     `json:"statusInfo"`
}

type EventsResp [][]interface{}

type Nothing = struct{}

type IdData struct {
	Id string `json:"id"`
}

type IdMsgdata struct {
	Id  string `json:"id"`
	Msg string `json:"msg"`
}
