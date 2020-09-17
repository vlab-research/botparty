package botparty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Event struct {
	Type string `json:"type"`
	Value *json.RawMessage `json:"value,omitempty"`
}

type ExternalEvent struct {
	User string `json:"user"`
	Page string `json:"page"`
	Event *Event `json:"event"`
}

func NewExternalEvent(user, page, type_ string, value *json.RawMessage) *ExternalEvent {
	return &ExternalEvent{user, page, &Event{type_, value}}
}

type BotParty struct {
	Client *http.Client 
	Botserver string
}

func NewBotParty(botserver string) *BotParty {
	return &BotParty{&http.Client{}, botserver}
}

func (eventer *BotParty) Send(e *ExternalEvent) error {
	body, err := json.Marshal(e)
	if err != nil {
		return err
	}

	resp, err := eventer.Client.Post(eventer.Botserver, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	code := resp.StatusCode
	if code != http.StatusOK {
		err := fmt.Errorf("Non 200 response from Botserver: %v", code)
		log.Print(err)
		return err
	}

	return nil
}
