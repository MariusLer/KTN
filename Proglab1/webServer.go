package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func handleRequest(conn net.Conn) {
	defer conn.Close() // lukker connection når funksjonen er ferdig
	fmt.Println("fikk en rq")
	// bufer to hold data
	buf := make([]byte, 1024)

	//reads the incomming connection into the buffer
	_, err := conn.Read(buf)

	if err != nil {
		fmt.Println("Eror reading:", err.Error())
		return
	}

	//leser ord nummer 2 fra bufferet
	filepath := strings.Split(string(buf[0:]), " ")[1]

	//filebuf := make([]byte, 1024)
	b, err := ioutil.ReadFile(filepath[1:])
	if err != nil {
		fmt.Println("eror reading file")
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	} else {
		fmt.Println("sender ting tilbake")
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		conn.Write(b)
		conn.Write([]byte("\r\n"))
	}
}

func server() {
	// Creates a tcp server socket
	ln, err := net.Listen("tcp", "192.168.38.105:8080") // change to use another ip if you want
	if err != nil {
		fmt.Println("feil")
		os.Exit(1)
	}

	// close socket when application closes
	defer ln.Close()
	for {
		fmt.Println("Ready to serve")

		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting")
			os.Exit(1)
		}
		//handles the request in a new thread
		go handleRequest(conn)
	}
}

func main() {
	server()
}
