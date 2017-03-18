package main

import (
	"bufio"
	"fmt"
	"messages"
	"net"
	"os"
	"strings"
	"time"
)

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
	splitInput = strings.Split(input, " ")
	var msg messages.ClientPayload
	msg.Req
	req_beg := strings.Index(input, "\\")
	if req_beg == 0 {
		req_end := strings.Index(input, " ")
		msg.Request = input[req_begin+1 : req_end]
		msg.Content = input[req_end+1:]
	} else {
		msg.Request = "msg"
		msg.Content = "input"
	}
}

func main() {
	server := "127.0.0.1:30000"

	conn := connectToServer(server)
	defer conn.Close()

	inputCh := make(chan string)

	go inputListener(inputCh)

	for {
		select {
		case input := <-inputCh: // Set up msg, json then send it
			msg := handleInput(input)

		}
	}
}
