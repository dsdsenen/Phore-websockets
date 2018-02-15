// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

// Addresses is a channel used to register an address to a websocket client
type Addresses struct {
	client  *Client
	address []string
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	subscribedToBlocks map[*Client]bool

	// Output messages to the clients.
	broadcastBlocks chan []byte

	// Register requests from the clients.
	registerBlock chan *Client

	// Unregister requests from clients.
	unregister     chan *Client
	unsubscribeAll chan *Client
}

func newHub() *Hub {
	// registerAddresses := make(chan Register)
	return &Hub{
		broadcastBlocks:    make(chan []byte),
		registerBlock:      make(chan *Client),
		unsubscribeAll:     make(chan *Client),
		subscribedToBlocks: make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.registerBlock:
			h.subscribedToBlocks[client] = true
		case client := <-h.unsubscribeAll:
			if _, ok := h.subscribedToBlocks[client]; ok {
				delete(h.subscribedToBlocks, client)
				close(client.send)
			}
		case message := <-h.broadcastBlocks:
			for client := range h.subscribedToBlocks {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.subscribedToBlocks, client)
				}
			}
		}
	}
}
