/*
A shell that connects to btcwallet using websockets.

Connects to a testnet btcwallet instance by default, add the -simnet flag for simnet,
-mainnet flag for mainnet.

*/

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/monetas/tools/btcwallet_websocket"
	"github.com/monetas/websocket"
)

func prettyPrintJSON(b []byte) error {
	var dat map[string]interface{}
	if err := json.Unmarshal(b, &dat); err != nil {
		return err
	}
	b, err := json.MarshalIndent(dat, "", "...  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func main() {
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

	conn, err := btcwallet_websocket.Connect(port)
	if err != nil {
		fmt.Println(err)
		return
	}

	discardReplies := true
	var discardRepliesMutex sync.Mutex
	replyChan := make(chan []byte)

	// A goroutine that loops forever reading replies and, if discardReplies is
	// false, sends them over replyChan.
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
			}
			discardRepliesMutex.Lock()
			if !discardReplies {
				replyChan <- msg
			}
			discardRepliesMutex.Unlock()
		}
	}()

	fmt.Println("Welcome to the btcwallet shell. Just enter your JSON requests here and "+
	            "we'll print all the replies that come within 500ms. Everything else is "+
				"discarded.")
	prompt := "btcwallet> "
	fmt.Print(prompt)
	// Loop reading lines, send them to btcwallet and print all the replies we get in 500ms,
	// after which we start discarding replies again.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		discardRepliesMutex.Lock()
		discardReplies = false
		discardRepliesMutex.Unlock()
		conn.WriteMessage(websocket.TextMessage, scanner.Bytes())
		timeout := time.After(500 * time.Millisecond)
		func() {
			for {
				select {
				case msg := <-replyChan:
					prettyPrintJSON(msg)
				case <-timeout:
					return
				}
			}
		}()
		discardRepliesMutex.Lock()
		discardReplies = true
		discardRepliesMutex.Unlock()
		fmt.Print("\n" + prompt)
	}
	if scanner.Err() != nil {
		fmt.Println(scanner.Err())
	}
}
