package websockets

import (
	"encoding/json"
	"log"

	"github.com/phoreproject/btcd/btcjson"
	"github.com/phoreproject/btcd/chaincfg/chainhash"
	"github.com/phoreproject/btcd/rpcclient"
)

// NotificationBlockHandler used to notify blocks
func NotificationBlockHandler(hub *Hub, client *rpcclient.Client, blockID string) {
	hash, err := chainhash.NewHashFromStr(blockID)
	if err != nil {
		log.Println("Error parsing the hash: ", err)
		return
	}

	data, err := client.GetBlockVerbose(hash)
	if err != nil {
		log.Println("Error getting the block: ", err)
		return
	}

	// Broadcast messages to subscribed clients asynchronously
	go broadcastBlocks(hub, data)
	go broadcastTransactions(hub, client, data)
}

// NotificationMempoolHandler used to notify mempool blocks
func NotificationMempoolHandler(hub *Hub, client *rpcclient.Client, txID string) {
	broadcastTransaction(hub, client, txID, true)
}

func broadcastBlocks(hub *Hub, data *btcjson.GetBlockVerboseResult) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Println("Error getting block info: ", err)
		return
	}
	hub.broadcastBlock <- []byte(string(jsonData))
}

func broadcastTransactions(hub *Hub, client *rpcclient.Client, data *btcjson.GetBlockVerboseResult) {
	for _, txID := range data.Tx {
		broadcastTransaction(hub, client, txID, false)
	}
}

func broadcastTransaction(hub *Hub, client *rpcclient.Client, txID string, memPool bool) {
	hashTx, err := chainhash.NewHashFromStr(txID)
	tx, err := client.GetRawTransactionVerbose(hashTx)
	if err != nil {
		log.Println("Error getting transaction: ", err)
		return
	}
	for _, transaction := range tx.Vout {
		for _, address := range transaction.ScriptPubKey.Addresses {
			jsonTx, _ := json.Marshal(tx)
			broadcastTransaction := BroadcastAddressMessage{address: address, message: []byte(string(jsonTx)), memPool: memPool}
			hub.broadcastAddress <- broadcastTransaction
			hub.broadcastBloom <- broadcastTransaction
		}
	}
}
