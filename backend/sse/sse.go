package sse
// Based on https://gist.github.com/ismasan/3fb75381cd2deb6bfa9c
// Copyright (c) 2017 Ismael Celis

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Broker struct {
	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan []byte

	// New client connections are pushed to this channel
	newClients chan chan []byte

	// Closed client connections are pushed to this channel
	closingClients chan chan []byte

	// Client connections registry
	clients map[chan []byte]bool
}

func NewServer() (broker *Broker) {
	// Instantiate a broker
	broker = &Broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()

	return
}

func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.newClients:

			// A new client has connected.
			// Register their message channel
			broker.clients[s] = true
			log.Printf("Client added. %d registered clients", len(broker.clients))
		case s := <-broker.closingClients:

			// A client has dettached and we want to
			// stop sending them messages.
			delete(broker.clients, s)
			log.Printf("Removed client. %d registered clients", len(broker.clients))
		case event := <-broker.Notifier:

			// We got a new event from the outside!
			// Send event to all connected clients
			for clientMessageChan := range broker.clients {
				clientMessageChan <- event
			}
		}
	}
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Transaction struct {
	Sender    int     `json:"sender"`
	Recipient int     `json:"recipient"`
	Amount    float64 `json:"amount"`
}

type RefreshEvent struct {
	AccountNr int `json:"account_nr"`
}

func (broker *Broker) Stream(w http.ResponseWriter, r *http.Request) {
	// Check if the ResponseWriter supports flushing.
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Each connection registers its own message channel with the Broker's connections registry
	messageChan := make(chan []byte)

	// Signal the broker that we have a new connection
	broker.newClients <- messageChan

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		broker.closingClients <- messageChan
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	for {
		select {
		// Listen to connection close and un-register messageChan
		case <-r.Context().Done():
			// remove this client from the map of connected clients
			broker.closingClients <- messageChan
			return

		// Listen for incoming messages from messageChan
		case msg := <-messageChan:
			// Write to the ResponseWriter
			// Server Sent Events compatible
			fmt.Fprintf(w, "data: %s\n\n", msg)
			// Flush the data immediatly instead of buffering it for later.
			flusher.Flush()
		}
	}
}

func (broker *Broker) BroadcastMessage(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send the message to the broker via Notifier channel
	j, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	broker.Notifier <- []byte(j)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Message sent"))
}

func (broker *Broker) BroadcastEvent(event string) {
	broker.Notifier <- []byte(event)
}

// To test the server, run the following commands in separate terminals:
// Start listening to the stream
//     $ curl -N http://localhost:<port>/stream
// Send a message
//     $ curl -X POST -H "Content-Type: application/json" -d '{"type": "transaction", "data": "{"sender": 1, "recipient": 2, "amount": 100"}"}' http://localhost:8000/messages

func (broker *Broker) SendMessage(eventType string, payload interface{}) {
	msg := Message{
		Type: eventType,
		Data: payload,
	}
	j, _ := json.Marshal(msg)
	broker.Notifier <- j
}
