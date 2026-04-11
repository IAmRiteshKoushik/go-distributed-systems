package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Node struct {
	NodeID    string
	NodeIDs   []string
	NextMsgID int
	mu        sync.Mutex
	outMu     sync.Mutex

	// For handling re-ordering of messages despite concurrency
	inbox chan Message
	wg    sync.WaitGroup
}

type Message struct {
	Src  string                 `json:"src"`
	Dest string                 `json:"dest"`
	Body map[string]interface{} `json:"body"`
}

func (n *Node) Send(dest string, body map[string]interface{}) {
	n.mu.Lock()
	body["msg_id"] = n.NextMsgID
	n.NextMsgID++
	n.mu.Unlock()

	msg := Message{Src: n.NodeID, Dest: dest, Body: body}
	output, _ := json.Marshal(msg)

	n.outMu.Lock()
	fmt.Println(string(output))
	n.outMu.Unlock()
}

func (n *Node) Reply(request Message, body map[string]interface{}) {
	if msgID, ok := request.Body["msg_id"].(float64); ok {
		body["in_reply_to"] = int(msgID)
	}
	n.Send(request.Src, body)
}

// Spin up workers and keep ready. All of these are listening to the a
// buffered channel. You can spin-up as many workers as you would like.
func (n *Node) worker() {
	defer n.wg.Done()

	for msg := range n.inbox {
		n.HandleMessage(msg)
	}
}

func (n *Node) HandleMessage(msg Message) {
	body := msg.Body

	msgType, ok := body["type"].(string)
	if !ok {
		return
	}

	switch msgType {
	case "init":
		// After the init message, setup all the connections
		n.mu.Lock()
		n.NodeID = body["node_id"].(string)
		rawIDs := body["node_ids"].([]interface{})
		for _, id := range rawIDs {
			n.NodeIDs = append(n.NodeIDs, id.(string))
		}
		n.mu.Unlock()

		// And then reply to it
		reply := map[string]interface{}{
			"type": "init_ok",
		}
		n.Reply(msg, reply)
	case "echo":
		reply := map[string]interface{}{
			"type": "echo_ok",
			"echo": msg.Body["echo"],
		}
		n.Reply(msg, reply)
	default:
	}
}

func main() {
	node := &Node{
		inbox: make(chan Message, 100),
	}

	numWorkers := 1
	node.wg.Add(numWorkers)

	for range numWorkers {
		go node.worker()
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			continue
		}
		node.inbox <- msg
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Scanner error:", err)
	}

	// Close the buffered channel to avoid memory leak and then wait till the
	// worker is done
	close(node.inbox)
	node.wg.Wait()
}
