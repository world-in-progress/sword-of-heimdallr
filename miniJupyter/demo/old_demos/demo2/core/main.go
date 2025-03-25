package main

import (
	"encoding/json"
	"log"

	"github.com/pebbe/zmq4"
)

type Core struct {
    routerSocket *zmq4.Socket
}

type Request struct {
    Type   string                 `json:"type"`
    Params map[string]interface{} `json:"params"`
}

type Response struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

func NewCore() (*Core, error) {
    router, err := zmq4.NewSocket(zmq4.ROUTER)
    if err != nil {
        return nil, err
    }

    err = router.Bind("tcp://*:5555")
    if err != nil {
        router.Close()
        return nil, err
    }

    return &Core{
        routerSocket: router,
    }, nil
}

func (c *Core) processRequest(message []byte) []byte {
    var request Request
    err := json.Unmarshal(message, &request)
    if err != nil {
        response := Response{
            Status:  "error",
            Message: "Invalid JSON format",
        }
        responseBytes, _ := json.Marshal(response)
        return responseBytes
    }

    // 处理请求
    response := Response{
        Status:  "success",
        Message: "Task executed successfully",
        Data: map[string]interface{}{
            "receivedParams": request.Params,
            "processedAt":   "some_timestamp",
        },
    }

    responseBytes, _ := json.Marshal(response)
    return responseBytes
}

func (c *Core) Start() {
    log.Println("Core service started on port 5555")
    
    for {
        // 接收消息（第一帧是身份标识）
        identity, err := c.routerSocket.RecvBytes(0)
        if err != nil {
            log.Printf("Error receiving identity: %v", err)
            continue
        }

        // 接收实际消息内容
        message, err := c.routerSocket.RecvBytes(0)
        if err != nil {
            log.Printf("Error receiving message: %v", err)
            continue
        }

        log.Printf("Received message: %s", string(message))

        // 处理请求
        result := c.processRequest(message)

        // 发送响应（需要包含原始身份标识）
        c.routerSocket.SendBytes(identity, zmq4.SNDMORE)
        c.routerSocket.SendBytes(result, 0)
    }
}

func main() {
    core, err := NewCore()
    if err != nil {
        log.Fatal(err)
    }

    core.Start()
}