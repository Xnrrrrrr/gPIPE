package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	// Define the source and destination ports
	sourcePort := "8080"
	destinationPort := "9090"

	// Start the server
	err := redirector(sourcePort, destinationPort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func redirector(sourcePort, destinationPort string) error {
	// Listen for incoming connections on the source port
	listener, err := net.Listen("tcp", ":"+sourcePort)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Printf("Redirecting from :%s to :%s\n", sourcePort, destinationPort)

	for {
		// Accept incoming connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting connection: %v\n", err)
			continue
		}

		// Print when a connection is accepted
		fmt.Println("Accepted connection.")

		// Handle the connection in a goroutine
		go handleConnection(conn, destinationPort)
	}
}

func handleConnection(conn net.Conn, destinationPort string) {
	defer conn.Close()

	// Print when handling a connection
	fmt.Println("Handling connection.")

	// Connect to the destination port
	destConn, err := net.Dial("tcp", ":"+destinationPort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to destination port: %v\n", err)
		return
	}
	defer destConn.Close()

	// Print when successfully connected to the destination port
	fmt.Println("Connected to destination port.")

	// Copy data bidirectionally
	go copyData(conn, destConn)
	copyData(destConn, conn)
}

func copyData(dst, src net.Conn) {
	_, err := io.Copy(dst, src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error copying data: %v\n", err)
	}
}
