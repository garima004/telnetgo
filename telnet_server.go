package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const (
	PORT     = ":8080"
	USERNAME = "admin"
	PASSWORD = "password"
)

// Object represents a simple structure with an ID and a Name
type Object struct {
	ID   int
	Name string
}

var objects = []Object{
	{ID: 1, Name: "preset 1"},
	{ID: 2, Name: "preset 2"},
	{ID: 3, Name: "preset 3"},
}

func main() {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Telnet server started on port", PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	conn.Write([]byte("Username: "))
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	conn.Write([]byte("Password: "))
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if username == USERNAME && password == PASSWORD {
		conn.Write([]byte("Login successful\n"))
		displayMenu(conn)
		for {
			conn.Write([]byte("Enter command: "))
			command, _ := reader.ReadString('\n')
			command = strings.TrimSpace(command)
			switch command {
			case "1":
				displayObjects(conn)
			case "2":
				displayPreset(conn)
			case "3":
				addObject(conn, reader)
			case "4":
				conn.Write([]byte("Goodbye!\n"))
				return
			default:
				conn.Write([]byte("Unknown command\n"))
				displayMenu(conn)
			}
		}
	} else {
		conn.Write([]byte("Login failed\n"))
	}
}

func displayMenu(conn net.Conn) {
	menu := "Menu:\n" +
		"1. show presetId\n" +
		"2. display preset\n" +
		"3. Add new preset\n" +
		"4. Exit\n"
	conn.Write([]byte(menu))
}

func displayObjects(conn net.Conn) {
	for _, obj := range objects {
		conn.Write([]byte(fmt.Sprintf("ID: %d, Name: %s\n", obj.ID, obj.Name)))
	}
}

func addObject(conn net.Conn, reader *bufio.Reader) {
	conn.Write([]byte("Enter new object name: "))
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		conn.Write([]byte("Invalid name. Returning to menu.\n"))
		return
	}

	id := len(objects) + 1
	objects = append(objects, Object{ID: id, Name: name})
	conn.Write([]byte(fmt.Sprintf("Object added: ID: %d, Name: %s\n", id, name)))
}

func displayPreset(conn net.Conn) {
	conn.Write([]byte("Prest successfully shown.\n"))
}
