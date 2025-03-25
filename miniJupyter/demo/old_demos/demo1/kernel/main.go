package main

import (
	"encoding/json"
	"fmt"

	zmq "github.com/pebbe/zmq4"
)

type Message struct {
    MsgType string      `json:"msg_type"`
    Content interface{} `json:"content"`
    MsgId   string      `json:"msg_id"`
}

type ExecuteRequest struct {
    Code   string `json:"code"`
    Silent bool   `json:"silent"`
}

func main() {
    // 创建ZMQ上下文
    router, _ := zmq.NewSocket(zmq.ROUTER)
    defer router.Close()
    router.Bind("tcp://*:5555")

    pub, _ := zmq.NewSocket(zmq.PUB)
    defer pub.Close()
    pub.Bind("tcp://*:5556")

    fmt.Println("Kernel started...")

    for {
        // 接收消息
        identity, _ := router.Recv(0)
        _, _ = router.Recv(0) // 空帧
        data, _ := router.Recv(0)

        var msg Message
        json.Unmarshal([]byte(data), &msg)

        // 处理消息
        switch msg.MsgType {
        case "heartbeat":
            // 响应心跳
            router.Send(identity, zmq.SNDMORE)
            router.Send("", zmq.SNDMORE)
            router.Send(data, 0)  // 直接返回收到的心跳消息
            
        case "execute_request":
            var execReq ExecuteRequest
            content, _ := json.Marshal(msg.Content)
            json.Unmarshal(content, &execReq)

            // 执行代码（这里简化处理）
            result := fmt.Sprintf("Executed: %s", execReq.Code)
            pub_result := fmt.Sprintf("PUB Executed: %s", execReq.Code)

            // 发送结果
            response := Message{
                MsgType: "execute_reply",
                Content: result,
                MsgId:   msg.MsgId,
            }
            responseData, _ := json.Marshal(response)

            pub_msg := Message{
                MsgType: "execute_reply",
                Content: pub_result,
                MsgId:   msg.MsgId,
            }
            pub_responseData, _ := json.Marshal(pub_msg)

            router.Send(identity, zmq.SNDMORE)
            router.Send("", zmq.SNDMORE)
            router.Send(string(responseData), 0)

            // 广播输出
            pub.Send(string(pub_responseData), 0)
        }
    }
}