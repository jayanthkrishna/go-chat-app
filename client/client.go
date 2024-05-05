package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/websocket"
)

type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}

type Client struct {
	ws          *websocket.Conn
	phonenumber string
}

func main() {
	origin := "http://localhost/"
	url := "ws://localhost:3000/ws"

	log.Println("Enter your phone number here")
	phonebuf, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	var phonenumber string = string(phonebuf)
	if err != nil {
		log.Println("erro while reading user phone number", err)
	}

	url = fmt.Sprintf("%s?phone=%s", url, phonenumber)

	ws, err := websocket.Dial(url, "", origin)

	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	fmt.Println("Connected to the server ")
	// done := make(chan struct{})

	client := Client{
		ws:          ws,
		phonenumber: phonenumber,
	}
	go client.ReadLoop()

	go client.WriteLoop()

}

func (c *Client) ReadLoop() {
	defer c.ws.Close()

	for {
		buf := make([]byte, 1024)
		n, err := c.ws.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				log.Println("Server closed the connection")
			} else {
				log.Println("Error receiving message:", err)
			}

			return
		}

		var message Message

		err = json.Unmarshal(buf[:n], &message)

		if err != nil {
			log.Println("Error while unmarshalling the message from server")
			continue
		}

		log.Printf("Received message from %s : %s\n", message.From, message.Message)

	}
}

func (c *Client) WriteLoop() {

	defer c.ws.Close()
	// defer close(done)
	for {
		reader := bufio.NewReader(os.Stdin)
		tobuf, err := reader.ReadBytes('\n')
		messagebuf, err := reader.ReadBytes('\n')

		message := Message{
			From:    c.phonenumber,
			To:      string(tobuf),
			Message: string(messagebuf),
		}

		msg, err := json.Marshal(&message)

		if err != nil {
			log.Println("Error while marshalling the message")
			continue
		}

		_, err = c.ws.Write(msg)
		if err != nil {
			if err.Error() == "EOF" {
				log.Println("Server closed the connection")
			} else {
				log.Println("Error writing to Server:", err)
			}

			return

		}
		log.Printf("Message written to %s", message.To)
	}

}
