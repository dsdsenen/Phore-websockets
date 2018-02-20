package main

import (
	"fmt"

	"github.com/phoreproject/btcd/wire"
	"github.com/phoreproject/btcutil/bloom"
)

// subscribebloom <FilterInHex> <HashFuncs> <Tweak> <Flags>
// but serializing it across websockets
// https://godoc.org/github.com/btcsuite/btcutil/bloom#Filter.MsgFilterLoad
// which returns this: https://godoc.org/github.com/btcsuite/btcd/wire#MsgFilterLoad
// Package wire
// Package wire implements the bitcoin wire protocol.
// which needs to be serialized
// with json/hex is probably fine
// then sent over the network through websockets
// and loaded with https://godoc.org/github.com/btcsuite/btcutil/bloom#LoadFilter
func subscribeBloom(client *Client, msg string) {
	fmt.Println(msg)
	// filter := strings.Split(msg, " ")
	// elems := filter[1]
	// hashFuncs := filter[2]
	// tweak := filter[3]
	// flags := filter[4]
	// flags =
	// bloomFilter := bloom.NewFilter(elems, 0, 0, flags)
	bloomFilter := bloom.NewFilter(100000000, 0, 0.01, wire.BloomUpdateNone)
	fmt.Println(bloomFilter)
	registerBloom := RegisterBloom{client, bloomFilter}
	client.hub.registerBloom <- registerBloom

}

// SubscribeAddress is used for a client to subscribe to any events happening to an address
func subscribeAddress(client *Client, addr string) {
	// fmt.Println("One new address registered", client, addr)
	register := RegisterAddress{client, addr}
	client.hub.registerAddress <- register
}

func subscribeBlock(client *Client) {
	// fmt.Println("One new client registered", client)
	client.hub.registerBlock <- client
}

func unsubscribeAll(client *Client) {
	client.hub.unsubscribeAll <- client
}
