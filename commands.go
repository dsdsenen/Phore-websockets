package main

import (
	"fmt"
)

func subscribeBloom(addr string) {
	fmt.Println(addr)
}

// SubscribeAddress is used for a client to subscribe to any events happening to an address
func subscribeAddress(addr string) {
	fmt.Println(addr)
}

func subscribeBlock(blockHash string) {
	fmt.Println(blockHash)
}

func unsubscribeAll() {

}
