/*
Websocket Client is a utility to connect to btcwallet using websockets

Connects to testnet by default, add the -simnet flag for simnet,
-mainnet flag for mainnet.

*/

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/monetas/tools/btcwallet_websocket"
	"github.com/monetas/websocket"
)

type T struct {
	Msg   string
	Count int
}

func main() {
	// message is the JSON to be sent to the websocket connection
	var simnet bool
	var mainnet bool
	var port int
	flag.BoolVar(&simnet, "simnet", false, "connect to simnet")
	flag.BoolVar(&mainnet, "mainnet", false, "connect to mainnet")
	flag.IntVar(&port, "port", 18332, "specific port to connect to")
	flag.Parse()

	if mainnet {
		port = 8332
	} else if simnet {
		port = 18554
	}

	arguments := flag.Args()
	if len(arguments) != 1 {
		fmt.Println("Usage: websocket <JSON to send to btcwallet websocket server>")
		return
	}
	message := []byte(arguments[0])

	conn, err := btcwallet_websocket.Connect(port)
	if err != nil {
		log.Fatal(err)
	}

	// send message to websocket connection.
	conn.WriteMessage(websocket.TextMessage, message)

	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			log.Fatal(err)
		}

		m := string(msg)

		log.Println(m)
	}
}
