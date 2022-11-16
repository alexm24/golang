package models

type Params struct {
	Channel string      `json:"channel"`
	Data    interface{} `json:"data"`
}

type Centrifugo struct {
	Method string `json:"method"`
	Params `json:"params"`
}

type ActionCentrifugo struct {
	Type    string      `json:"type,omitempty"`
	Payload interface{} `json:"payload"`
}
