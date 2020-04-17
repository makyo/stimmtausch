package eventstream

type Event struct {
	Type      string      `json:"type"`
	ClientID  string      `json:"client_id"`
	Timestamp string      `json:"client_id"`
	Payload   interface{} `json:"payload"`
}

type SignalPayload struct {
	Name    string   `json:"name"`
	Payload []string `json:"payload"`
}

type MessagePayload interface{}
