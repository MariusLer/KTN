package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
	"unicode"

	"github.com/mariusler/KTN/Project/msg"
)

type incommingMsg struct {
	Conn net.Conn
	Msg  messages.ClientPayload
}

func main() {
	clients := make(map[net.Conn]string)
	//msgHistory := make([]messages.ServerPayload)

	newConnCh := make(chan net.Conn)
	inCommingMsgCh := make(chan incommingMsg)
	removeDeadClientCh := make(chan net.Conn)

	go connListener(newConnCh)

	for {
		select {
		case newConn := <-newConnCh:
			clients[newConn] = ""
			go clientListener(inCommingMsgCh, newConn, removeDeadClientCh)
			fmt.Println("New client")
			fmt.Println(clients)
		case deadConn := <-removeDeadClientCh:
			fmt.Println("User :", clients[deadConn], " Logged off")
			delete(clients, deadConn)
		case msgReceived := <-inCommingMsgCh:
			conn := msgReceived.Conn
			msg := msgReceived.Msg
			switch msg.Request {
			case "login":
				if isValidUserName(msg.Content) {
					if clients[conn] == "" { // ikke lurt her, gjør på nytt
						clients[conn] = msg.Content
						sendMessage(messages.ServerPayload{time.Now().String(), "Server", "info", []string{"Login succesful"}}, conn)
					} else {
						sendMessage(messages.ServerPayload{time.Now().String(), "Server", "error", []string{"Already logged in"}}, conn)
					}
					if isNameTaken(clients, msg.Content) {
						sendMessage(messages.ServerPayload{time.Now().String(), "Server", "error", []string{"Username taken"}}, conn)
					}
				} else {
					sendMessage(messages.ServerPayload{time.Now().String(), "Server", "error", []string{"Username not valid"}}, conn)
				}
			case "logout":
				if clients[conn] == "" {
					sendMessage(messages.ServerPayload{time.Now().String(), "Server", "error", []string{"Not logged in"}}, conn)
				} else {
					conn.Close()
				}
			case "names":
				names := getNameList(clients)
				sendMessage(messages.ServerPayload{time.Now().String(), "Server", "info", []string{names}}, conn)
			case "help":
			case "msg":
			}
		}
	}
}

func sendMessage(msg messages.ServerPayload, conn net.Conn) {
	bytes := serverPayloadToNetworkMsg(msg)
	conn.Write(bytes)
}

func getNameList(clients map[net.Conn]string) string {
	var names string
	for _, name := range clients {
		names += name + "\n"
	}
	return names
}

func isNameTaken(clients map[net.Conn]string, userName string) bool {
	for _, name := range clients {
		if name == userName {
			return false
		}
	}
	return true
}

func serverPayloadToNetworkMsg(msg messages.ServerPayload) []byte {
	bytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Error json Marshal", err)
		return []byte("")
	}
	return bytes
}

func connListener(newConnCh chan<- net.Conn) {
	ln, err := net.Listen("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("Error")
		os.Exit(1)
	}
	defer ln.Close()
	fmt.Println("Ready to listen to connections")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting")
		} else {
			newConnCh <- conn
		}
	}
}
func isValidUserName(name string) bool {
	if name == "" {
		return false
	}
	for _, letter := range name {
		if (letter < 'a' || letter > 'z') && (letter < 'A' || letter > 'Z') { // it is not a letter
			if !unicode.IsDigit(letter) {
				return false
			}
		}
	}
	return true
}

func clientListener(inCommingMsgCh chan<- incommingMsg, conn net.Conn, removeDeadClientCh chan<- net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 2048)
	var msg messages.ClientPayload

	for {
		bytes, err := conn.Read(buffer)

		if err != nil {
			removeDeadClientCh <- conn
			return
		}
		error := json.Unmarshal(buffer[:bytes], &msg)
		if error != nil {
			fmt.Println("Error Unmarshall", err)
			continue
		}
		inCommingMsgCh <- incommingMsg{conn, msg}
	}
}
