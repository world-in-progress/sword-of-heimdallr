package main

import (
	"log"
	"time"

	"miniJupyter/zmq/base"
	"miniJupyter/zmq/mode"
)

func main() {
    // 加载配置
    config, err := base.LoadConfig("zmq/config/config.yaml")
    if err != nil {
        log.Fatal(err)
    }

    // 启动各个组件
    // go mode.RunRouter(config.Zmq.RouterAddress)
    go mode.RunHeartbeatServer(config.Zmq.RouterAddress)
    
    // 等待一段时间让服务器启动
    time.Sleep(time.Second)
    
    // go mode.RunDealer(config.Zmq.DealerAddress)
    go mode.RunHeartbeatClient(config.Zmq.DealerAddress, time.Duration(config.Zmq.HeartbeatInterval) * time.Second)

    // 防止程序退出
    select {}
}