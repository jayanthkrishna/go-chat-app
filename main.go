package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	mrand "math/rand"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	connections map[string]map[string]*websocket.Conn
}

type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}

func NewServer() *Server {
	return &Server{
		connections: make(map[string]map[string]*websocket.Conn),
	}
}
func (s *Server) handleWS(ws *websocket.Conn) {

	fmt.Println("new incoming connection from client :", ws.RemoteAddr())
	defer ws.Close()
	phonenumber := ws.Request().URL.Query().Get("phone")

	if phonenumber == "" {
		log.Println("Phone number is empty")
		return
	}

	log.Printf("Phone number trying to connect :%s\n", phonenumber)
	otp := generateOTP()

	sendOTP(phonenumber, otp)

	var OTP struct {
		OTP string `json:"otp"`
	}

	msg := make([]byte, 100)
	n, _ := ws.Read(msg)

	err := json.Unmarshal(msg[:n], &OTP)

	if err != nil {
		log.Println("Error while unmarshalling")
	}

	if otp != OTP.OTP {
		log.Println("Invalid OTP")
		return
	}

	token := generateToken()

	if s.connections[phonenumber] == nil {
		s.connections[phonenumber] = make(map[string]*websocket.Conn)
	}

	s.connections[phonenumber][token] = ws

	var Token = map[string]string{"token": token}
	jsonToken, _ := json.Marshal(Token)

	if _, err := ws.Write([]byte(jsonToken)); err != nil {
		log.Println("failed while sending token to user")
		return
	}

	s.Listen(ws)

}

func generateOTP() string {
	return fmt.Sprintf("%06d", mrand.Intn(1000000))
}

// sendOTP simulates sending the OTP to the user's phone number (for demonstration purposes)
func sendOTP(phoneNumber, otp string) {
	log.Printf("Sending OTP %s to phone number %s\n", otp, phoneNumber)
	// Simulate sending OTP via SMS or another communication channel
}

func generateToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("Error generating token:", err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (s *Server) Listen(ws *websocket.Conn) {
	defer ws.Close()

	for {

		var message Message

		buf := make([]byte, 1024)
		n, err := ws.Read(buf)

		if err != nil {

			if err.Error() == "EOF" {
				log.Println("Client closed the connection")
			} else {
				log.Println("Error receiving message:", err)
			}

			return
		}

		err = json.Unmarshal(buf[:n], &message)

		if err != nil {
			log.Println("Error while Unmarshalling the message")
		}

		if s.connections[message.To] != nil {
			for _, conn := range s.connections[message.To] {
				if _, err := conn.Write(buf[:n]); err != nil {
					log.Println("error while sending the data to receiver")

				}

			}
		} else {
			log.Println("receiver not found")
		}
		fmt.Printf("messsage received from %s : %s\n", message.From, message.Message)

	}

}

func main() {

	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))

	fmt.Println("Server listening at port 3000")
	http.ListenAndServe(":3000", nil)

}
