package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
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
	msgHistory := make([]messages.ServerPayload, 0)

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
			fmt.Println("Map of clients", clients)
		case deadConn := <-removeDeadClientCh:
			fmt.Println("User :", clients[deadConn], " Logged off")
			delete(clients, deadConn)
		case msgReceived := <-inCommingMsgCh:
			conn := msgReceived.Conn
			msg := msgReceived.Msg
			switch msg.Request {
			case "login":
				if isValidUserName(msg.Content) {
					if clients[conn] == "" {
						if isNameTaken(clients, msg.Content) {
							sendMessage(messages.ServerPayload{Timestamp: time.Now().String(), Sender: "Server", Response: "error", Content: "Username taken"}, conn)
						} else {
							clients[conn] = msg.Content
							sendMessage(messages.ServerPayload{Timestamp: time.Now().String(), Sender: "Server", Response: "info", Content: "Login succesful"}, conn)
							sendChatHistory(msgHistory, conn)
						}
					} else {
						sendMessage(messages.ServerPayload{Timestamp: time.Now().String(), Sender: "Server", Response: "error", Content: "Already logged in"}, conn)
					}

				} else {
					sendMessage(messages.ServerPayload{Timestamp: time.Now().String(), Sender: "Server", Response: "error", Content: "Username not valid"}, conn)
				}
			case "logout":
				if clients[conn] == "" {
					sendMessage(messages.ServerPayload{Timestamp: time.Now().String(), Sender: "Server", Response: "error", Content: "Not logged in"}, conn)
				} else {
					conn.Close()
				}
			case "names":
				names := getNameList(clients)
				sendMessage(messages.ServerPayload{Timestamp: time.Now().String(), Sender: "Server", Response: "info", Content: names}, conn)
			case "help":
				sendMessage(messages.ServerPayload{Timestamp: time.Now().String(), Sender: "Server", Response: "info", Content: "The commands supported are login <username> logout names help and msg <message>, use backslash in front of the commands"}, conn)
			case "msg":
				response := messages.ServerPayload{Timestamp: time.Now().String(), Sender: clients[conn], Response: "message", Content: msg.Content}
				msgHistory = append(msgHistory, response)
				broadcastMsg(clients, response)
			default:
				sendMessage(messages.ServerPayload{Timestamp: time.Now().String(), Sender: "Server", Response: "error", Content: "Unknown command"}, conn)
			}
		}
	}
}
func broadcastMsg(clients map[net.Conn]string, msg messages.ServerPayload) {
	for conn := range clients {
		sendMessage(msg, conn)
	}
}

func sendChatHistory(msgHistory []messages.ServerPayload, conn net.Conn) {
	var list = make([][]byte, len(msgHistory))
	for index, msg := range msgHistory {
		bytes, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("Error json marshal history", err)
			return
		}
		list[index] = bytes
	}
	msg := messages.HistoryPayload{Timestamp: time.Now().String(), Sender: "Server", Response: "History", Content: list}
	bytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Error json marshal msg", err)
		return
	}
	conn.Write(bytes)
}

func sendMessage(msg messages.ServerPayload, conn net.Conn) {
	bytes := serverPayloadToNetworkMsg(msg)
	conn.Write(bytes)
}

func serverPayloadToNetworkMsg(msg messages.ServerPayload) []byte {
	bytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Error json Marshal", err)
		return []byte("")
	}
	return bytes
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
			return true
		}
	}
	return false
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
	if name == "" || strings.ToLower(name) == "server" {
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
