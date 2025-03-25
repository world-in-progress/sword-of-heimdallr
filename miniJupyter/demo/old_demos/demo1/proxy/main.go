package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	zmq "github.com/pebbe/zmq4"
)

type Message struct {
	MsgType string      `json:"msg_type"`
	Content interface{} `json:"content"`
	MsgId   string      `json:"msg_id"`
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients    = make(map[*Client]bool)
	clientsMux sync.Mutex
)

func (c *Client) write() {
	defer func() {
		c.conn.Close()
		clientsMux.Lock()
		delete(clients, c)
		clientsMux.Unlock()
	}()

	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
}

// 添加心跳检测
func (c *Client) heartbeat(dealer *zmq.Socket) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 发送心跳请求
		dealer.Send("", zmq.SNDMORE)
		dealer.Send(`{"msg_type": "heartbeat"}`, 0)
		
		// 接收响应
		_, err := dealer.Recv(0)
		if err != nil {
			log.Println("Heartbeat failed:", err)
			c.send <- []byte(`{"msg_type": "heartbeat_status", "content": "failed"}`)
			continue
		}
		response, err := dealer.Recv(0)
		if err != nil {
			log.Println("Heartbeat failed:", err)
			c.send <- []byte(`{"msg_type": "heartbeat_status", "content": "failed"}`)
			continue
		}
		
		c.send <- []byte(`{"msg_type": "heartbeat_status", "content": "alive"}`)
		log.Printf("Heartbeat success: %s\n", response)
	}
}

func main() {
	// 连接到Kernel
	dealer, _ := zmq.NewSocket(zmq.DEALER)
	defer dealer.Close()
	dealer.Connect("tcp://localhost:5555")

	sub, _ := zmq.NewSocket(zmq.SUB)
	defer sub.Close()
	sub.Connect("tcp://localhost:5556")
	sub.SetSubscribe("")

	// 处理从Kernel接收的消息
	go func() {
		for {
			msg, _ := sub.Recv(0)
			broadcastToClients([]byte(msg))
		}
	}()

	// WebSocket处理
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		client := &Client{
			conn: conn,
			send: make(chan []byte, 256),
		}

		clientsMux.Lock()
		clients[client] = true
		clientsMux.Unlock()

		// 启动写入协程
		go client.write()
		// 启动心跳检测
		go client.heartbeat(dealer)

		// 读取循环
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			// 转发到Kernel
			dealer.Send("", zmq.SNDMORE)
			dealer.Send(string(message), 0)

			// 接收Kernel响应
			_, _ = dealer.Recv(0)
			response, _ := dealer.Recv(0)

			// 发送到client的channel
			client.send <- []byte(response)
		}

		// 清理
		close(client.send)
	})

	fmt.Println("Proxy started on :8080")
	http.ListenAndServe(":8080", nil)
}

func broadcastToClients(message []byte) {
	clientsMux.Lock()
	defer clientsMux.Unlock()
	for client := range clients {
		select {
		case client.send <- message:
		default:
			// 如果channel已满，关闭连接
			close(client.send)
			delete(clients, client)
		}
	}
}