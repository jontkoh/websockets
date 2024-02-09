package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
)

func readMessage(conn net.Conn) (string, error) {
	header := make([]byte, 2)
	_, err := conn.Read(header)
	if err != nil {
		return "", err
	}
	var maskKey [4]byte
	var payloadLength int

	maskBit := header[1] & 0x80
	length := int(header[1] & 0x7F)

	// Adjust for extended payload lengths
	if length == 126 {
		extended := make([]byte, 2)
		_, err = conn.Read(extended)
		if err != nil {
			return "", err
		}
		payloadLength = int(binary.BigEndian.Uint16(extended))
		_, err = conn.Read(maskKey[:]) // Read the masking key after the 2 bytes
	} else if length == 127 {
		extended := make([]byte, 8)
		_, err = conn.Read(extended)
		if err != nil {
			return "", err
		}
		payloadLength = int(binary.BigEndian.Uint64(extended))
		_, err = conn.Read(maskKey[:]) // Read the masking key after the 8 bytes
	} else {
		payloadLength = length
		_, err = conn.Read(maskKey[:]) // Read the masking key directly if length <= 125
	}

	payload := make([]byte, length)
	_, err = conn.Read(payload)
	if err != nil {
		return "", err
	}

	// Unmask the payload if it's masked
	if maskBit != 0 && payloadLength > 0 {
		for i := range payload {
			payload[i] ^= maskKey[i%4]
		}
	}

	return string(payload), nil
}

func writeMessage(conn net.Conn, message string) error {
	length := len(message)
	var header []byte
	if length <= 125 {
		header = []byte{0x81, byte(length)}
	} else if length <= 65535 {
		header = []byte{0x81, 126, byte(length >> 8), byte(length & 0xFF)}
	} else {
		header = []byte{0x81, 127, 0, 0, 0, 0, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length & 0xFF)}
	}

	_, err := conn.Write(append(header, []byte(message)...))
	return err
}

func initiateWebSocket(w http.ResponseWriter, r *http.Request) {
	// ALL required for protocol switching
	upgradeHeader := r.Header.Get("Upgrade")
	connectionUpgradeHeader := r.Header.Get(("Connection"))
	webSocketKey := r.Header.Get("Sec-WebSocket-Key")

	if upgradeHeader != "websocket" || webSocketKey == "" || connectionUpgradeHeader == "" {
		http.Error(w, "Missing headers", http.StatusInternalServerError)
		return
	}
	hash := sha1.New()

	// PUBLIC GUID
	hash.Write([]byte(webSocketKey + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))

	acceptKey := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	w.Header().Set("Upgrade", "websocket")
	w.Header().Set("Connection", "Upgrade")
	w.Header().Set("Sec-WebSocket-Accept", acceptKey)
	w.WriteHeader(http.StatusSwitchingProtocols)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "http.Hijacker interface is not supported", http.StatusInternalServerError)
		return
	}

	conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go func() {
		defer conn.Close()

		for {
			// Example of reading a message
			addClient(&conn)
			message, err := readMessage(conn)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				return
			}
			fmt.Printf("Received message: %s\n", message)

			// Example of sending a message back
			broadcastMessage(&conn, message)
			if err := writeMessage(conn, message); err != nil {
				log.Printf("Error sending message: %v", err)
				return
			}
		}
	}()
}

type NewMessageData struct {
	Message string `json:"message"`
}

var (
	// Protects access to the clients map
	clientsMutex sync.Mutex
	// Maps to keep track of active connections
	clients = make(map[*net.Conn]bool)
)

// Add a new client connection
func addClient(conn *net.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	clients[conn] = true
}

// Remove a client connection
func removeClient(conn *net.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	delete(clients, conn)
}

// Broadcast a message to all clients (except the sender)
func broadcastMessage(sender *net.Conn, message string) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	for conn := range clients {
		if conn != sender { // Check if the current connection is not the sender
			writeMessage(*conn, message)
		}
	}
}
