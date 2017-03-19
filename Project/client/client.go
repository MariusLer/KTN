package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/mariusler/KTN/Project/msg"
)

func main() {
	server := "127.0.0.1:30000"

	conn := connectToServer(server)
	defer conn.Close()

	inputCh := make(chan string)
	incommingMsgCh := make(chan messages.ServerPayload)
	incommingHistoryCh := make(chan messages.HistoryPayload)
	closedConnCh := make(chan bool)

	go inputListener(inputCh)
	go messageListener(incommingMsgCh, incommingHistoryCh, conn, closedConnCh)

	for {
		select {
		case input := <-inputCh:
			msg := handleInput(input)
			bytes, err := json.Marshal(msg)
			if err != nil {
				fmt.Println("Error json marshall", err)
			}
			conn.Write(bytes)
		case msg := <-incommingMsgCh:
			printMsg(msg)
		case hMsg := <-incommingHistoryCh:
			var oldMessage messages.ServerPayload
			if len(hMsg.Content) != 0 {
				fmt.Print("[", hMsg.Timestamp[11:16], "] ", "<", "Server", "> ", hMsg.Response, " : ", "Chat history received, listing it below", "\n")
			}
			for _, byteObject := range hMsg.Content {
				err := json.Unmarshal(byteObject, &oldMessage)
				if err != nil {
					fmt.Println("Error unmarhsalling history")
					fmt.Println(byteObject)
					continue
				}
				printMsg(oldMessage)
			}
		case <-closedConnCh:
			return
		}
	}
}

func printMsg(msg messages.ServerPayload) {
	if msg.Sender != "Server" {
		fmt.Print("[", msg.Timestamp[11:16], "] ", "<", msg.Sender, "> ", msg.Content, "\n")
	} else {
		fmt.Print("[", msg.Timestamp[11:16], "] ", "<", "Server", "> ", msg.Response, " : ", msg.Content, "\n")
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

func messageListener(incommingMsgCh chan<- messages.ServerPayload, incommingHistoryCh chan<- messages.HistoryPayload, conn net.Conn, closedConnCh chan<- bool) {
	buffer := make([]byte, 1048576) // huge buffer 1 MiB
	var msg messages.ServerPayload
	var historyMsg messages.HistoryPayload
	for {
		bytes, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Closed connection/logged out")
			closedConnCh <- true
			return
		}
		// This part is messy. There are probably much better solutions
		errr := json.Unmarshal(buffer[:bytes], &msg)
		if errr != nil {
			json.Unmarshal(buffer[:bytes], &historyMsg)
			incommingHistoryCh <- historyMsg
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
		if reqEnd == -1 {
			msg.Request = input[reqBeg+1:]
		} else {
			msg.Request = input[reqBeg+1 : reqEnd]
		}
		if msg.Request == "login" {
			if len(splitInput) != 1 {
				msg.Content = input[reqEnd+1:]
			}
		}
	} else {
		msg.Request = "msg"
		msg.Content = input
	}
	return msg
}
