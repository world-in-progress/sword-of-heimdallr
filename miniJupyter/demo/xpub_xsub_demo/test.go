package main

import (
	"fmt"
	"log"
	"miniJupyter/zmq/mode"
	"time"
)

func main() {
    // 1. 创建 XPUB 节点
    pub, err := mode.NewXPublisher("tcp://*:5555")
    if err != nil {
        log.Fatal(err)
    }
    defer pub.Close()

    // 2. 设置不同主题的权限
    pub.SetTopicPermission("topic1", []string{"user1", "user2"})         // topic1 只允许 user1 和 user2 访问
    pub.SetTopicPermission("topic2", []string{"user2", "user3"})         // topic2 只允许 user2 和 user3 访问
    pub.SetTopicPermission("private", []string{"admin"})                 // private 主题只允许 admin 访问
    // topic3 不设置权限，默认所有用户可访问

    // 3. 启动 XPUB 节点的主循环
    go func() {
        if err := pub.Run(); err != nil {
            log.Printf("XPUB error: %v\n", err)
        }
    }()

    // 4. 创建多个 XSUB 节点
    sub1, err := mode.NewXSubscriber("tcp://localhost:5555", "user1")
    if err != nil {
        log.Fatal(err)
    }
    defer sub1.Close()

    sub2, err := mode.NewXSubscriber("tcp://localhost:5555", "user2")
    if err != nil {
        log.Fatal(err)
    }
    defer sub2.Close()

    sub3, err := mode.NewXSubscriber("tcp://localhost:5555", "user3")
    if err != nil {
        log.Fatal(err)
    }
    defer sub3.Close()

    adminSub, err := mode.NewXSubscriber("tcp://localhost:5555", "admin")
    if err != nil {
        log.Fatal(err)
    }
    defer adminSub.Close()

    // 5. 测试订阅场景
    fmt.Println("=== 测试订阅权限 ===")
    
    // user1 订阅测试
    testSubscribe(sub1, "topic1")  // 应该成功
    testSubscribe(sub1, "topic2")  // 应该失败
    testSubscribe(sub1, "topic3")  // 应该成功（无权限限制）
    
    // user2 订阅测试
    testSubscribe(sub2, "topic1")  // 应该成功
    testSubscribe(sub2, "topic2")  // 应该成功
    
    // user3 订阅测试
    testSubscribe(sub3, "topic1")  // 应该失败
    testSubscribe(sub3, "topic2")  // 应该成功
    
    // admin 订阅测试
    testSubscribe(adminSub, "private")  // 应该成功
    
    // 6. 等待订阅建立
    time.Sleep(time.Second)

    // 7. 发布消息测试
    fmt.Println("\n=== 测试消息发布 ===")
    testPublish(pub, "topic1", "Message for topic1")
    testPublish(pub, "topic2", "Message for topic2")
    testPublish(pub, "topic3", "Message for everyone")
    testPublish(pub, "private", "Secret message")

    // 8. 测试取消订阅
    fmt.Println("\n=== 测试取消订阅 ===")
    testUnsubscribe(sub1, "topic1")
    testUnsubscribe(sub2, "topic2")

    // 9. 测试动态修改权限
    fmt.Println("\n=== 测试动态修改权限 ===")
    pub.RemoveTopicPermission("topic1")  // 移除 topic1 的权限限制
    fmt.Println("Removed permissions for topic1")
    testSubscribe(sub3, "topic1")  // 现在应该可以成功

    // 10. 添加新的权限
    pub.SetTopicPermission("topic1", []string{"user3"})  // 只允许 user3 访问
    fmt.Println("Set new permissions for topic1")
    testSubscribe(sub1, "topic1")  // 现在应该失败

    // 等待一段时间以观察结果
    time.Sleep(time.Second * 2)
}

func testSubscribe(sub *mode.XSubscriberNode, topic string) {
    err := sub.Subscribe(topic)
    if err != nil {
        fmt.Printf("订阅失败 - UserID: %s, Topic: %s, Error: %v\n", "sub.userID", topic, err)
    } else {
        fmt.Printf("订阅成功 - UserID: %s, Topic: %s\n", "sub.userID", topic)
    }
}

func testUnsubscribe(sub *mode.XSubscriberNode, topic string) {
    err := sub.Unsubscribe(topic)
    if err != nil {
        fmt.Printf("取消订阅失败 - UserID: %s, Topic: %s, Error: %v\n", "sub.userID", topic, err)
    } else {
        fmt.Printf("取消订阅成功 - UserID: %s, Topic: %s\n", "sub.userID", topic)
    }
}

func testPublish(pub *mode.XPublisherNode, topic, message string) {
    err := pub.Publish(topic, message)
    if err != nil {
        fmt.Printf("发布失败 - Topic: %s, Message: %s, Error: %v\n", topic, message, err)
    } else {
        fmt.Printf("发布成功 - Topic: %s, Message: %s\n", topic, message)
    }
}