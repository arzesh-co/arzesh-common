package messaging

import (
	"encoding/json"
	errors2 "errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
)

func Send(message interface{}, stream, subject string) error {
	// Connect to NATS
	natsUrl := os.Getenv("NATSConnection")
	if natsUrl == "" {
		return errors2.New("connection config is not set")
	}
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return err
	}
	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return err
	}
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     stream,
		Subjects: []string{stream + ".*"},
	})
	if err != nil {
		return err
	}

	// Encode the person object as JSON
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Publish the message to the stream
	msg := &nats.Msg{
		Subject: subject,
		Data:    payload,
	}
	_, err = js.PublishMsg(msg)
	if err != nil {
		return err
	}

	return nil
}

func Consume(stream, subject string, action func(msg *nats.Msg)) {
	// Connect to NATS
	natsUrl := os.Getenv("NATSConnection")
	if natsUrl == "" {
		fmt.Println("connection config is not set")
		return
	}
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		fmt.Println("error to connection :", err.Error())
		return
	}

	// Create JetStream Context
	js, _ := nc.JetStream(nats.PublishAsyncMaxPending(256))

	// Get last sequence number
	_, err = js.StreamInfo(stream)
	if err != nil {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:        stream,
			AllowDirect: true,
			Subjects:    []string{stream + ".*"},
		})
		if err != nil {
			fmt.Println("error to sync 2 :", err.Error())
			return
		}
		_, err = js.StreamInfo(stream)
		if err != nil {
			fmt.Println("error to get stream info :", err.Error())
			return
		}
	}
	// Simple Async Ephemeral Consumer with StartSeq option
	_, err = js.Subscribe(subject, action, nats.Durable(stream)) // start from the next sequence number
	if err != nil {
		fmt.Println("error to subscribe :", err.Error())
		return
	}
	// Wait for messages
	fmt.Println("Waiting for messages...")
	select {}
}

type Histories struct {
	ClientToken string `json:"client_token"`
	UserToken   string `json:"user_token"`
	DocId       string `json:"doc_uuid"`
	ServiceKey  string `json:"service_key"`
	EntityName  string `json:"entity_name"`
	Action      string `json:"action"`
	VersionNo   int64  `json:"version_no"`
	Document    any    `json:"document"`
}

func SendHistory(histories *Histories) error {
	// Connect to NATS
	natsUrl := os.Getenv("NATSConnection")
	if natsUrl == "" {
		return errors2.New("connection config is not set")
	}
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return err
	}
	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return err
	}
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "histories",
		Subjects: []string{"histories.*"},
	})
	if err != nil {
		return err
	}

	// Encode the person object as JSON
	payload, err := json.Marshal(histories)
	if err != nil {
		return err
	}

	// Publish the message to the stream
	msg := &nats.Msg{
		Subject: "histories." + histories.Action,
		Data:    payload,
	}
	_, err = js.PublishMsg(msg)
	if err != nil {
		return err
	}

	return nil
}
