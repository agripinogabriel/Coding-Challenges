package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var basePath = "./www%s"
var basePage = "/index.html"
var defaultPort = "8080"
var defaultNetwork = "tcp"

func main() {

	port := readServerPort()

	// Listen for incoming connections on port 8080
	ln, err := net.Listen(defaultNetwork, port)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Accept incoming connections and handle them
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Close the connection when we're done
	defer conn.Close()

	// Read the incoming data
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		conn.Write([]byte("HTTP/1.1 500 Server Error\r\n"))
		return
	}

	// Print the incoming data
	fmt.Printf("Received: %s", buf)

	rawRequest := string(buf)
	firstLine := strings.Split(rawRequest, "\n")
	requestElements := strings.Split(firstLine[0], " ")
	requestedPath := requestElements[1]

	page, err := readPage(requestedPath)
	if err != nil {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n"))
		return
	}

	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	conn.Write([]byte("\r\n"))
	conn.Write(page)
}

func readPage(path string) (content []byte, err error) {
	if strings.Contains(path, basePage) == false {
		path = fmt.Sprintf("%s%s", path, basePage)
	}

	dat, err := os.ReadFile(fmt.Sprintf(basePath, path))
	if err != nil {
		return nil, err
	}

	return dat, nil
}

func readServerPort() string {
	args := os.Args[1:]

	if len(args) > 0 {
		return fmt.Sprintf(":%s", args[0])
	}

	return fmt.Sprintf(":%s", defaultPort)
}
