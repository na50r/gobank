package sse

// Original based on: https://github.com/plutov/packagemain/tree/master/30-sse
// YouTube video: https://www.youtube.com/watch?v=nvijc5J-JAQ

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Broker struct {
	cnt            int
	connected      chan ClientChannel
	disconnected   chan int
	clientChannels map[int]chan []byte
}

type ClientChannel struct {
	ID      int
	Channel chan []byte
}

type Message struct {
	Data interface{} `json:"data"`
}

func NewBroker() (broker *Broker) {
	broker = &Broker{
		connected:      make(chan ClientChannel),
		disconnected:   make(chan int),
		clientChannels: make(map[int]chan []byte),
		cnt:            0,
	}
	go broker.listen()
	return
}

// Goroutine manages multiple channels in the background
func (b *Broker) listen() {
	for {
		select {
		case cc := <-b.connected:
			b.clientChannels[cc.ID] = cc.Channel
			fmt.Printf("client connected (id=%d), total=%d\n", cc.ID, len(b.clientChannels))
		case remove := <-b.disconnected:
			delete(b.clientChannels, remove)
			fmt.Printf("client disconnected (id=%d), total=%d\n", remove, len(b.clientChannels))
		}
	}
}

func (b *Broker) createChannel(id int) ClientChannel {
	if id == 0 {
		id = b.cnt
	}
	ch := make(chan []byte)
	cc := ClientChannel{ID: id, Channel: ch}
	return cc
}

func (b *Broker) SSEHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	cc := b.createChannel(0)
	b.connected <- cc

	defer func() {
		b.disconnected <- cc.ID
	}()

	clientGone := r.Context().Done()

	rc := http.NewResponseController(w)

	for {
		select {
		case <-clientGone:
			return
		case data := <-cc.Channel:
			if _, err := fmt.Fprintf(w, "event:msg\ndata:%s\n\n", data); err != nil {
				log.Printf("unable to write: %s", err.Error())
				return
			}
			rc.Flush()
		}
	}
}

func (b *Broker) Publish(msg Message) {
	data, err := json.Marshal(msg.Data)
	if err != nil {
		log.Printf("unable to marshal: %s", err.Error())
		return
	}
	// Publish to all channels
	// NOTE: Not concurrent
	for _, channel := range b.clientChannels {
		channel <- data
	}
}

func (b *Broker) PublishEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var m Message
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	b.Publish(m)
	w.Write([]byte("Msg sent\n"))
}
