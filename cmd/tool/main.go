package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func connectAndSend(wg *sync.WaitGroup, serverURL string, messages int) {
	defer wg.Done()

	u, err := url.Parse(serverURL)
	if err != nil {
		log.Fatal("Error parsing URL:", err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("Error connecting to WebSocket server:", err)
		return
	}
	defer conn.Close()

	// Create a 1 KB message
	message := strings.Repeat("A", 1024) // A simple 1 KB message of repeated 'A' characters

	for i := 0; i < messages; i++ {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Write error:", err)
			return
		}

		_, _, err = conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: benchmark <ws://your-websocket-server/ws> <number_of_connections>")
		return
	}

	serverURL := os.Args[1]
	connections := 1500
	var wg sync.WaitGroup

	// Set the number of messages per connection to 15,000
	messagesPerConnection := 15000

	start := time.Now()
	for i := 0; i < connections; i++ {
		wg.Add(1)
		go connectAndSend(&wg, serverURL, messagesPerConnection)
	}

	wg.Wait()
	fmt.Printf("Benchmark completed in %v seconds.\n", time.Since(start).Seconds())
}
