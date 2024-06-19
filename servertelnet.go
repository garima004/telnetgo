package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/reiver/go-telnet"
)

// Constants for login credentials
const (
	port     = "8080"
	httpPort = "8081"
	username = "admin"
	password = "password"
)

// structure to represent a preset object
type preset struct {
	PresetId string `json:"presetId"`
	Pname    string `json:"pname"`
}

var presets = []preset{
	{"P1", "Preset1"},
	{"P2", "Preset2"},
	{"P3", "Preset3"},
}

// mutex on preset
//var presetsMutex sync.Mutex

// TelnetHandler struct to implement the telnet.Handler interface
type TelnetHandler struct{}

// ServeTELNET method to handle Telnet connections
func (handler TelnetHandler) ServeTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {
	if !authenticate(w, r) {
		return
	}
	showMenu(w, r)
}

// authenticate function to handle user authentication
func authenticate(w telnet.Writer, r telnet.Reader) bool {
	w.Write([]byte("Username: "))
	inputUsername, err := readLine(r)
	if err != nil {
		fmt.Println("Error reading username:", err)
		return false
	}
	inputUsername = strings.TrimSpace(inputUsername)
	fmt.Println("Received username:", inputUsername) // Debugging line
	w.Write([]byte("Password: "))
	inputPassword, err := readLine(r)
	if err != nil {
		fmt.Println("Error reading password:", err)
		return false
	}
	inputPassword = strings.TrimSpace(inputPassword)
	fmt.Println("Received password:", inputPassword) // Debugging line
	if inputUsername == username && inputPassword == password {
		w.Write([]byte("Welcome!\n"))
		return true
	}
	w.Write([]byte("Access denied!\n"))
	return false
}

// showMenu function to display the menu and handle user input
func showMenu(w telnet.Writer, r telnet.Reader) {
	for {
		w.Write([]byte("\nMenu:"))
		w.Write([]byte("\n1. List Presets"))
		w.Write([]byte("\n2. Run Presets"))
		w.Write([]byte("\n3. Add Presets"))
		w.Write([]byte("\n4. Exit"))
		w.Write([]byte("\nEnter your choice: "))

		choice, err := readLine(r)
		if err != nil {
			fmt.Println("Error reading choice:", err)
			return
		}
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			listPresets(w)
		case "2":
			runPreset(w, r)
		case "3":
			addPreset(w, r)
		case "4":
			w.Write([]byte("Goodbye!\n"))
			return
		default:
			w.Write([]byte("Invalid choice, please try again.\n"))
		}
	}
}

// listBooks function to display the list of books
func listPresets(w telnet.Writer) {
	//presetsMutex.Lock()
	//defer presetsMutex.Unlock()
	w.Write([]byte("\nList of Presets:\n"))
	for i, preset := range presets {
		w.Write([]byte(fmt.Sprintf("%d.  %s  %s\n", i+1, preset.PresetId, preset.Pname)))
	}
}

func runPreset(w telnet.Writer, r telnet.Reader) {
	w.Write([]byte("Enter preset id: "))
	id, err := readLine(r)
	if err != nil {
		fmt.Println("Error reading preset id:", err)
		w.Write([]byte("Failed to run preset.\n"))
		return
	}
	id = strings.TrimSpace(id)
	w.Write([]byte("Preset ran successfully " + id))

}

func addPreset(w telnet.Writer, r telnet.Reader) {
	w.Write([]byte("Enter preset id: "))
	id, err := readLine(r)
	if err != nil {
		fmt.Println("Error reading preset id:", err)
		w.Write([]byte("Failed to add preset.\n"))
		return
	}
	id = strings.TrimSpace(id)

	w.Write([]byte("Enter preset name: "))
	name, err := readLine(r)
	if err != nil {
		fmt.Println("Error reading preset name:", err)
		w.Write([]byte("Failed to add preset.\n"))
		return
	}
	name = strings.TrimSpace(name)
	//presetsMutex.Lock()
	presets = append(presets, preset{PresetId: id, Pname: name})
	//presetsMutex.Unlock()
	w.Write([]byte("Preset added successfully!\n"))
}

// readLine function to read a line of input from the Telnet client
func readLine(r telnet.Reader) (string, error) {
	var line []byte
	var buffer [1]byte
	for {
		_, err := r.Read(buffer[:])
		if err != nil {
			return "", err
		}
		if buffer[0] == '\n' {
			break
		}
		line = append(line, buffer[0])
	}
	return string(line), nil
}

func listPresetHTTP(w http.ResponseWriter, r *http.Request) {
	//presetsMutex.Lock()
	//defer presetsMutex.Unlock()
	fmt.Println("inside request")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(presets)
}

// main function to start the Telnet server
func main() {
	go func() {
		fmt.Println("Telnet server listening on port", port)
		err := telnet.ListenAndServe(":"+port, TelnetHandler{})
		if err != nil {
			fmt.Println("Error starting Telnet server:", err)
		}
	}()

	fmt.Println("Http server listneing")
	// Start HTTP server
	http.HandleFunc("/presets", listPresetHTTP)
	fmt.Println("HTTP server listening on port 8081")
	err1 := http.ListenAndServe(":8081", nil)
	if err1 != nil {
		fmt.Println("Error starting HTTP server:", err1)
	} else {
		fmt.Println("Http serevr started")
	}

}
