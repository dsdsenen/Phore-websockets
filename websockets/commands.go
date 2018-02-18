package websockets

import (
	"encoding/hex"
	"errors"
	"strconv"

	"github.com/phoreproject/btcd/wire"
	"github.com/phoreproject/btcutil/bloom"
)

func subscribeBloom(client *Client, args []string) error {
	// Syntax: subscribeBloom <filterHex> <HashFuncs> <Tweak>

	if len(args) != 3 {
		return errors.New("Incorrect number of arguments")
	}

	var bloomBytes []byte
	bloomBytes, err := hex.DecodeString(args[0])
	if err != nil {
		return errors.New("Could not decode bloom filter")
	}

	hashFuncsInt, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.New("Could not parse HashFuncs")
	}
	hashFuncs := uint32(hashFuncsInt)

	tweakInt, err := strconv.Atoi(args[2])
	if err != nil {
		return errors.New("Could not parse Tweak")
	}
	tweak := uint32(tweakInt)

	filter := bloom.LoadFilter(&wire.MsgFilterLoad{
		Filter:    bloomBytes,
		HashFuncs: hashFuncs,
		Tweak:     tweak,
		Flags:     wire.BloomUpdateNone,
	})

	client.hub.registerBloom <- RegisterBloom{client: client, bloom: filter}
	return nil
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
