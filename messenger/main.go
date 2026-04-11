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
}

type Message struct {
	Src  string                 `json:"src"`
	Dest string                 `json:"dest"`
	Body map[string]interface{} `json:"body"`
}

type Body struct {
	Type      string                 `json:"type"`
	MsgID     interface{}            `json:"msg_id,omitempty"`
	InReplyTo interface{}            `json:"in_reply_to,omitempty"`
	Extra     map[string]interface{} `json:"-"`
}

func (b Body) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["type"] = b.Type
	if b.InReplyTo != nil {
		m["in_reply_to"] = b.InReplyTo
	}
	m["msg_id"] = b.MsgID
	for k, v := range b.Extra {
		m[k] = v
	}
	return json.Marshal(m)
}

func (n *Node) Send(dest string, body map[string]interface{}) {
	n.mu.Lock()
	defer n.mu.Unlock()

	extra := make(map[string]interface{})
	for k, v := range body {
		if k != "type" {
			extra[k] = v
		}
	}

	msg := struct {
		Src  string `json:"src"`
		Dest string `json:"dest"`
		Body Body   `json:"body"`
	}{
		Src:  n.NodeID,
		Dest: dest,
		Body: Body{
			Type:  body["type"].(string),
			MsgID: n.NextMsgID,
			Extra: extra,
		},
	}
	if inReplyTo, ok := body["in_reply_to"]; ok {
		msg.Body.InReplyTo = inReplyTo
	}

	b, err := json.Marshal(msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "marshal error:", err)
		return
	}
	fmt.Println(string(b))
	n.NextMsgID++
}

func (n *Node) Reply(request Message, body map[string]interface{}) {
	body["in_reply_to"] = request.Body["msg_id"]
	n.Send(request.Src, body)
}

func main() {
	node := &Node{}
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			continue
		}

		msgType, _ := msg.Body["type"].(string)
		if msgType == "init" {
			node.NodeID = msg.Body["node_id"].(string)
			ids := msg.Body["node_ids"].([]interface{})
			for _, id := range ids {
				node.NodeIDs = append(node.NodeIDs, id.(string))
			}

			resp := make(map[string]interface{})
			resp["type"] = "init_ok"
			node.Reply(msg, resp)
		}
	}
}
