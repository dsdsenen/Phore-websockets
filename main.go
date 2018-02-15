// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/phoreproject/btcd/rpcclient"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	var connCfg = &rpcclient.ConnConfig{
		// Phore RPC Daemon
		Host: "127.0.0.1:11772",
		// Phore RPC Proxy
		// Host:                 "rpc.phore.io/rpc",
		HTTPPostMode:         true,
		User:                 "phorerpc",
		Pass:                 "JCiM652B1gW1bbbxLHwdnpETFNs3HoGndUGS2Ef2J8jq",
		DisableTLS:           true,
		DisableAutoReconnect: false,
		DisableConnectOnNew:  false,
	}

	client, _ := rpcclient.New(connCfg, nil)

	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r, client)
	})
	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		for k := range r.Form {
			notificationBlockHandler(hub, client, k)
			return
		}
	})

	log.Println("Starting Websockets Server...")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
