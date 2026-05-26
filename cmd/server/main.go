package main

import (
	"bufio"
	"fmt"
	"github.com/ouchbouk/tinydb"
	"net"
	"strings"
)

type Server struct {
	db *tinydb.DB
}

func main() {
	db, err := tinydb.Open("tiny.db")
	s := &Server{
		db: db,
	}
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer s.db.Close()
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

		value, err := s.db.Get([]byte(parts[1]))
		if err != nil {
			return "ERR " + err.Error()
		}

		return "OK " + string(value)

	case "SET":
		if len(parts) != 3 {
			return "ERR usage: SET key value"
		}

		err := s.db.Set([]byte(parts[1]), []byte(parts[2]))
		if err != nil {
			return "ERR " + err.Error()
		}

		return "OK"

	case "DELETE":
		if len(parts) != 2 {
			return "ERR usage: DELETE key value"
		}

		err := s.db.Delete([]byte(parts[1]))
		if err != nil {
			return "ERR " + err.Error()
		}
		return "OK"

	default:
		return "ERR unknown command"
	}
}
