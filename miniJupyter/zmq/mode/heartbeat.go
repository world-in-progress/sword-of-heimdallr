package mode

import (
	"fmt"
	"miniJupyter/zmq/base"
	"time"

	zmq "github.com/pebbe/zmq4"
)

// RunHeartbeatServer 心跳检测服务器
func RunHeartbeatServer(address string) {
	publisher, err := base.NewZmqNode(zmq.PUB, address, true)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()

	for {
		err = publisher.Send("heartbeat")
		if err != nil {
			fmt.Printf("Error publishing heartbeat: %v\n", err)
		}
		time.Sleep(time.Second) // 每秒发送一次心跳
	}
}

// RunHeartbeatClient 心跳检测客户端
func RunHeartbeatClient(address string, timeout time.Duration) chan bool {
	subscriber, err := base.NewZmqNode(zmq.SUB, address, false)
	if err != nil {
		panic(err)
	}
	defer subscriber.Close()

	err = subscriber.SetSubscribe("heartbeat")
	if err != nil {
		panic(err)
	}

	// 创建用于通知断连状态的channel
	disconnectedChan := make(chan bool)
	
	lastHeartbeat := time.Now()
	go func() {
		for {
			_, err := subscriber.Receive()
			if err != nil {
				fmt.Printf("Error receiving heartbeat: %v\n", err)
				continue
			}
			lastHeartbeat = time.Now()
		}
	}()

	// 监控心跳超时
	go func() {
		for {
			if time.Since(lastHeartbeat) > timeout {
				fmt.Println("Heartbeat timeout detected")
				disconnectedChan <- true  // 发送断连信号
				// 这里可以添加断开连接或重连的逻辑
			}
			time.Sleep(timeout / 2)
		}
	}()

	return disconnectedChan
}
