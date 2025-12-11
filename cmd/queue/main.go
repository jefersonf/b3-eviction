package queue

import (
	"b3e/internal/queue"
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {

	// 1. Initialize the core engine
	q := queue.NewMemoryQueue()

	// 2. Start TCP listener
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("Voting Queue running on :9000")

	for {
		// 3. Accept Connections
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		// 4. Handle concurrently
		go handleConnection(conn, q)
	}
}

func handleConnection(conn net.Conn, q *queue.MemoryQueue) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		text := scanner.Text()
		parts := strings.SplitN(text, " ", 2)
		command := strings.ToUpper(parts[0])

		switch command {
		case "PUSH":
			if len(parts) < 2 {
				fmt.Fprintln(conn, "ERR missing payload")
				continue
			}
			q.Enqueue(parts[1])
			fmt.Fprintln(conn, "OK")
		case "POP":
			val, ok := q.Dequeue()
			if !ok {
				fmt.Fprintln(conn, "NIL")
			} else {
				fmt.Fprintln(conn, val)
			}
		default:
			fmt.Fprintln(conn, "ERR unknown command")
		}
	}
}
