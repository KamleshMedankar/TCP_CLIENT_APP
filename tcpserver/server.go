package tcpserver

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
)

func StartMultipleTCPServers() {
	var wg sync.WaitGroup
	for port := 8000; port < 8010; port++ {
		if port == 8005 {
			continue
		}
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			startServer(p)
		}(port)
	}
	wg.Wait()
}

func startServer(port int) {
	address := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Failed to start server on port %d: %v\n", port, err)
		return
	}
	fmt.Printf("[Server %d] Listening on %s\n", port, address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("[Server %d] Accept error: %v\n", port, err)
			continue
		}
		go handleConnection(conn, port)
	}
}

func handleConnection(conn net.Conn, port int) {
	defer func() {
		fmt.Printf("[Server %d] Closing connection\n", port)
		conn.Close()
	}()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		if err != io.EOF {
			fmt.Printf("[Server %d] Read error: %v\n", port, err)
		}
		return
	}

	message := string(buf[:n])
	fmt.Printf("[Server %d] Received: %s\n", port, message)

	// Respond with ACK
	ack := fmt.Sprintf("ACK from port %d\n", port)
	_, err = conn.Write([]byte(ack))
	if err != nil {
		fmt.Printf("[Server %d] Write error: %v\n", port, err)
	}
}
