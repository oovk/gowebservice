package product

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

type message struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

func productSocket(ws *websocket.Conn) {

	done := make(chan struct{}) //to check if the connection is closed!
	fmt.Println("New websocket connection established")
	//listen for incoming data on the websockets
	go func(c *websocket.Conn) { //websocket returns a EOF error whwn client closes the connection
		for {
			var msg message
			if err := websocket.JSON.Receive(ws, &msg); err != nil {
				log.Print(err)
				break
			}
			fmt.Printf("Received Message %s\n", msg.Data)
		}
		close(done)
	}(ws)
	//continously get top 10 rows from the data after each 10s
loop:
	for {
		select {
		case <-done:
			fmt.Print("connection closed, break out of loop")
			break loop
		default:
			products, err := GetTopTenProducts()
			if err != nil {
				log.Print(err)
				break
			}
			if err := websocket.JSON.Send(ws, products); err != nil {
				log.Print(err)
				break
			}
			time.Sleep(10 * time.Second) //sleep for 10s and then refresh the data again
		}

	}
	fmt.Print("closing the websocket")
	defer ws.Close() //cleans up the websocket connection
}
