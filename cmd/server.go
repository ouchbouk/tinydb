package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Server struct {
	mu   sync.RWMutex
	data map[string]string
}

func main() {

	s := &Server{
		data: make(map[string]string),
	}
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	defer ln.Close()
	fmt.Println("server listening on :8080")

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println("accept error: ", err)
			continue
		}

		go s.handleConn(conn)
	}

}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("command:", line)

		response := s.handleCommand(line)
		conn.Write([]byte(response + "\n"))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
	}
}

func (s *Server) handleCommand(line string) string {
	parts := strings.Fields(line)

	if len(parts) == 0 {
		return "ERR empty command"
	}

	switch parts[0] {
	case "PING":
		return "PONG"

	case "GET":
		if len(parts) != 2 {
			return "ERR usage: GET key"
		}

		s.mu.RLock()
		value, ok := s.data[parts[1]]
		s.mu.RUnlock()

		if !ok {
			return "ERR not found"
		}

		return "OK " + value

	case "SET":
		if len(parts) != 3 {
			return "ERR usage: SET key value"
		}

		s.mu.Lock()
		s.data[parts[1]] = parts[2]
		s.mu.Unlock()

		return "OK"

	default:
		return "ERR unknown command"
	}
}
