package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// Message represents a Maelstrom message
type Message struct {
	Src  string                 `json:"src"`
	Dest string                 `json:"dest"`
	Body map[string]interface{} `json:"body"`
}

func main() {
	// TODO: Read JSON messages from stdin
	// Each line is a complete JSON message
	// Parse and print: PARSED: src|dest|body_type
	// Log details to stderr for debugging

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		var msg Message
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
			continue
		}

		// Extract the message body
		msgBody := "unknown"
		if msgBodyType, ok := msg.Body["type"]; ok {
			if confirmStr, ok := msgBodyType.(string); ok {
				msgBody = confirmStr
			}
		}
		msgBody, ok := msg.Body["type"].(string)
		if !ok {
			msgBody = "unknown"
		}

		fmt.Printf("PARSED: %s|%s|%v\n", msg.Src, msg.Dest, msgBody)
	}
	// For handling EOF situations
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Scanner error:", err)
	}
}
