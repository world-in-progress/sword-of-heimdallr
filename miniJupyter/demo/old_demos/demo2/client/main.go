package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pebbe/zmq4"
)

type Client struct {
    dealerSocket *zmq4.Socket
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // 允许所有跨域请求，生产环境需要更严格的检查
    },
}

func NewClient() (*Client, error) {
    dealer, err := zmq4.NewSocket(zmq4.DEALER)
    if err != nil {
        return nil, err
    }
    
    err = dealer.Connect("tcp://localhost:5555")
    if err != nil {
        dealer.Close()
        return nil, err
    }

    return &Client{
        dealerSocket: dealer,
    }, nil
}

func (c *Client) handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket upgrade error: %v", err)
        return
    }
    defer conn.Close()

    log.Printf("New WebSocket connection established")

    for {
        // 读取WebSocket消息
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Printf("WebSocket read error: %v", err)
            break
        }

        // 发送到Core服务
        log.Printf("Sending message to Core: %s", string(message))
        _, err = c.dealerSocket.SendBytes(message, 0)
        if err != nil {
            log.Printf("ZMQ send error: %v", err)
            continue
        }

        // 接收Core服务的响应
        response, err := c.dealerSocket.RecvBytes(0)
        if err != nil {
            log.Printf("ZMQ receive error: %v", err)
            continue
        }

        // 发送响应回WebSocket客户端
        err = conn.WriteMessage(websocket.TextMessage, response)
        if err != nil {
            log.Printf("WebSocket write error: %v", err)
            break
        }
    }
}

func main() {
    client, err := NewClient()
    if err != nil {
        log.Fatal(err)
    }

    // 静态文件服务
    http.Handle("/", http.FileServer(http.Dir("../frontend")))
    
    // WebSocket端点
    http.HandleFunc("/ws", client.handleWebSocket)

    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}