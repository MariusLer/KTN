package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"./msg"
)

func main() {
	server := "127.0.0.1:30000"

	conn := connectToServer(server)
	defer conn.Close()

	inputCh := make(chan string)
	incommingMsgCh := make(chan messages.ServerPayload)

	go inputListener(inputCh)
	go messageListener(incommingMsgCh, conn)

	for {
		select {
		case input := <-inputCh: // Set up msg, json then send it
			msg := handleInput(input)
			bytes, err := json.Marshal(msg)
			if err != nil {
				fmt.Println("Error json marshall", err)
			}
			conn.Write(bytes)
		case msg := <-incommingMsgCh:
			fmt.Println("msg received yo", msg)
		}
	}
}

func inputListener(inputCh chan<- string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error :", err)
			continue
		}
		input = strings.Replace(input, "\n", "", -1) // will work on unix and windows
		inputCh <- input
	}
}

func messageListener(incommingMsgCh chan<- messages.ServerPayload, conn net.Conn) {
	buffer := make([]byte, 2048)
	var msg messages.ServerPayload
	for {
		bytes, err := conn.Read(buffer)
		if err != nil {
			conn.Close()
			return
		}
		errr := json.Unmarshal(buffer[:bytes], msg)
		if errr != nil {
			fmt.Println("Error Unmarshall", err)
			continue
		}
		incommingMsgCh <- msg
	}
}

func connectToServer(ipAndPort string) net.Conn {
	for {
		conn, err := net.Dial("tcp", ipAndPort)
		if err == nil {
			return conn
		}
		fmt.Println("Error connecting to server")
		time.Sleep(time.Second)
	}
}

func handleInput(input string) messages.ClientPayload {
	splitInput := strings.Split(input, " ")
	var msg messages.ClientPayload
	reqBeg := strings.Index(input, "\\")
	if reqBeg == 0 {
		reqEnd := strings.Index(input, " ")
		msg.Request = input[reqBeg+1 : reqEnd]
		if len(splitInput) != 1 {
			msg.Content = input[reqEnd+1:]
		}
	} else {
		msg.Request = "msg"
		msg.Content = "input"
	}
	return msg
}
