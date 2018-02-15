package main

import (
	"fmt"
)

func subscribeBloom(client *Client, addr string) {
	fmt.Println(addr)
}

// SubscribeAddress is used for a client to subscribe to any events happening to an address
func subscribeAddress(client *Client, addr string) {
	fmt.Println("One new address registered", client, addr)
	// client.hub.registerAddress <-
}

func subscribeBlock(client *Client) {
	fmt.Println("One new client registered", client)
	client.hub.registerBlock <- client
}

func unsubscribeAll(client *Client) {
	client.hub.unsubscribeAll <- client
}
