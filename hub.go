// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

// RegisterAddress is a channel used to register an address to a websocket client
type RegisterAddress struct {
	client  *Client
	address string
}

// BroadcastAddressMessage used to receive message of addresses
type BroadcastAddressMessage struct {
	address string
	message []byte
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	subscribedToBlocks  map[*Client]bool
	subscribedToAddress map[string][]*Client

	// Output messages to the clients.
	broadcastBlock   chan []byte
	broadcastAddress chan BroadcastAddressMessage

	// Register requests from the clients.
	registerBlock   chan *Client
	registerAddress chan RegisterAddress

	// Unregister requests from clients.
	unregister     chan *Client
	unsubscribeAll chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcastBlock:      make(chan []byte),
		broadcastAddress:    make(chan BroadcastAddressMessage),
		registerBlock:       make(chan *Client),
		registerAddress:     make(chan RegisterAddress),
		unsubscribeAll:      make(chan *Client),
		subscribedToBlocks:  make(map[*Client]bool),
		subscribedToAddress: make(map[string][]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.registerBlock:
			h.subscribedToBlocks[client] = true
		case registerAddress := <-h.registerAddress:
			addr := registerAddress.address
			h.subscribedToAddress[addr] = append(h.subscribedToAddress[addr], registerAddress.client)
		case client := <-h.unsubscribeAll:
			if _, ok := h.subscribedToBlocks[client]; ok {
				delete(h.subscribedToBlocks, client)
				close(client.send)
			}
		case message := <-h.broadcastBlock:
			for client := range h.subscribedToBlocks {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.subscribedToBlocks, client)
				}
			}
		case broadcastAddress := <-h.broadcastAddress:
			addr := broadcastAddress.address
			for _, client := range h.subscribedToAddress[addr] {
				select {
				case client.send <- broadcastAddress.message:
				default:
					deleteClientFromAddress(client, addr)
					close(client.send)

				}
			}
		case client := <-h.unsubscribeAll:
			delete(h.subscribedToBlocks, client)
			// TODO: Improve this delete method
			for address, clients := range h.subscribedToAddress {
				if clientInSlice(client, clients) {
					deleteClientFromAddress(client, address)
				}
			}
		}
	}
}

func deleteClientFromAddress(client *Client, addr string) {
	var i int
	for j, v := range client.hub.subscribedToAddress[addr] {
		if v == client {
			i = j
		}
	}
	client.hub.subscribedToAddress[addr] = append(client.hub.subscribedToAddress[addr][:i], client.hub.subscribedToAddress[addr][i+1:]...)
}

func clientInSlice(client *Client, list []*Client) bool {
	for _, b := range list {
		if b == client {
			return true
		}
	}
	return false
}
