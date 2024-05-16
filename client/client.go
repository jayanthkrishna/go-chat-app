package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

type OTP struct {
	OTP string `json:"otp"`
}

func main() {

	log.Println("Enter your phone number here")
	phonebuf, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	var phonenumber string = string(phonebuf[:len(phonebuf)-1])
	if err != nil {
		log.Println("error while reading user phone number", err)
	}

	url := "http://localhost:3001/gen-token?phone="

	token, err := GetToken(url, phonenumber)

	if err != nil {
		log.Println("Error while getting token frok auth server")
		return
	}

	origin := "http://localhost/"
	url = "ws://localhost:3000/ws"

	url = fmt.Sprintf("%s?phone=%s&token=%s", url, phonenumber, token)
	log.Println("URL :", url)

	ws, err := websocket.Dial(url, "", origin)

	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	// log.Println("Enter the OTP here.")
	// otpbuf, _ := bufio.NewReader(os.Stdin).ReadBytes('\n')

	// otpbuf = otpbuf[:len(otpbuf)-1]
	// otp := string(otpbuf)

	// log.Println("OTP: ", otp)

	// OTP := OTP{
	// 	OTP: otp,
	// }

	// otpJSON, _ := json.Marshal(&OTP)
	// ws.Write(otpJSON)
	fmt.Println("Connected to the server ")
	done := make(chan struct{})

	client := Client{
		ws:          ws,
		phonenumber: phonenumber,
	}
	go client.ReadLoop(done)

	go client.WriteLoop(done)

	<-done

}

func GetToken(url string, phonenumber string) (string, error) {

	url = url + phonenumber

	// Send GET request to the server
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	type Token struct {
		Token string `json:"string"`
	}
	var token Token
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		fmt.Println("Error decoding the response to json", err)
		return "", err
	}

	// Print the response body
	return token.Token, nil
}

func (c *Client) ReadLoop(done chan struct{}) {
	defer close(done)
	buf := make([]byte, 1024)
	for {

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

func (c *Client) WriteLoop(done chan struct{}) {

	defer close(done)
	reader := bufio.NewReader(os.Stdin)
	for {

		tobuf, _ := reader.ReadBytes('\n')
		tobuf = bytes.TrimRight(tobuf, "\n")
		messagebuf, _ := reader.ReadBytes('\n')
		messagebuf = bytes.TrimRight(messagebuf, "\n")

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
		log.Printf("Message written to %s : %s", message.To, message.Message)
	}

}
