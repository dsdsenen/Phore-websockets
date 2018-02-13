package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/phoreproject/btcd/chaincfg/chainhash"
	"github.com/phoreproject/btcd/rpcclient"
)

func notificationBlockHandler(hub *Hub, client *rpcclient.Client, blockID string) {
	hash, err := chainhash.NewHashFromStr(blockID)
	fmt.Println("HUB: %+v", hub.addresses)
	if err != nil {
		log.Println("Error parsing the hash: ", err)
		return
	}

	data, err := client.GetBlockVerbose(hash)
	if err != nil {
		log.Println("Error getting block: ", err)
		return
	}
	fmt.Printf("BLOCK INFO: %+v\n", data)

	for _, txID := range data.Tx {
		hash2, err := chainhash.NewHashFromStr(txID)
		log.Println("HASH: ", hash2)
		data2, err := client.GetRawTransactionVerbose(hash2)
		if err != nil {
			log.Println("Error getting transaction: ", err)
			return
		}
		for _, transaction := range data2.Vout {
			fmt.Printf("TRANSACTION: %+v\n", transaction)
			for _, address := range transaction.ScriptPubKey.Addresses {
				fmt.Println("ADDRESS", address)
				jsonTx, _ := json.Marshal(data2)
				fmt.Println("BROADCAST KRAI")
				hub.broadcast <- []byte(string(jsonTx))
			}
		}
	}

	// address, err := btcutil.DecodeAddress("PDqJskowZHNxufyWhL2aTHho72RHEDHKti", defaultNet)
	// fmt.Println(address, err)
	// // addr, err := btcutil.DecodeAddress(addrString, defaultNet)
	// if err != nil {
	// 	fmt.Println("Error")
	// }
	// fmt.Println(addr.String())
	// data, err := client.SearchRawTransactionsVerbose(address, 0, 10, false, false, make([]string, 0))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("Transaction 0 confirmations: ", data[0].Confirmations)
}
