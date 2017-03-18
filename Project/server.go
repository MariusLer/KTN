package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"./msg"
)

func main() {
	clients := make(map[net.Conn]string)
	//msgHistory := make([]messages.ServerPayload)

	newConnCh := make(chan net.Conn)
	inCommingMsgCh := make(chan messages.ClientPayload)

	go connListener(newConnCh)

	for {
		select {
		case newConn := <-newConnCh:
			clients[newConn] = ""
			go clientListener(inCommingMsgCh, newConn)
			fmt.Println("New client")
			fmt.Println(clients)
		case msg := <-inCommingMsgCh:
			switch msg.Request {
			case "login":
			case "logout":
			case "names":
			case "help":
			default:

			}
		}
	}
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

func clientListener(inCommingMsgCh chan<- messages.ClientPayload, conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 2048)
	var msg messages.ClientPayload

	for {
		bytes, err := conn.Read(buffer)

		if err != nil {
			fmt.Println("Error receiving, closing connection")
			return
		}
		error := json.Unmarshal(buffer[:bytes], msg)
		if error != nil {
			fmt.Println("Error Unmarshall", err)
			continue
		}
		inCommingMsgCh <- msg
	}
}
