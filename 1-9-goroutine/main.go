package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 并发和并行的概念
// 并发：多个任务在同一时间段内交替执行（单核CPU）
// 并行：多个任务在同一时间点同时执行（多核CPU）

// CSP 理论
// 不要通过共享内存来通信，而通过通信来共享内存
// 组成：goroutine（协程）+ channel（通道）

// goroutine 的试用场景
// 高并发的web服务：每个请求一个 goroutine 处理
// 实时消息推送：每个消息一个 goroutine 处理
// 数据处理管道：每个数据处理步骤一个 goroutine
// 定时任务调度：每个定时任务一个 goroutine

// channel 的使用
// 创建 channel
// 无缓冲 channel: 同步通信，发送和接收必须配对
// ch1 := make(chan type)
// 有缓冲 channel: 异步通信，缓冲区满的时候阻塞，发送和接收可以不配对
// ch2 := make(chan type, buffer)

// 操作
// 发送
// ch1 <- 42
// ch2 <- 42
// 接收
// value := <-ch1
// value := <-ch2

// 检查 channel 是否关闭
// value, ok := <-ch1

// select 多路复用
// select {
// case <-ch1:
//     // ch1 有数据
// case <-ch2:
//     // ch2 有数据
// default:
//     // 都没有数据
// }

type Order struct {
	ID       int
	UserID   string
	Amount   int
	Status   string
	CreateAt time.Time
}

// 生产者 生产订单
func orderProduct(orderChan chan<- Order, number int) {
	for i := 0; i < number; i++ {
		order := Order{
			ID:       i,
			UserID:   fmt.Sprintf("userId_%d", rand.Intn(100)),
			Amount:   rand.Intn(10),
			Status:   "pending",
			CreateAt: time.Now(),
		}
		orderChan <- order
	}
	close(orderChan)
}

func processOrder(orderChan <-chan Order, resultChan chan<- Order) {
	for order := range orderChan {
		order.Status = "processing"
		fmt.Printf("处理订单: ID=%d, 用户ID=%s, 金额=%d\n", order.ID, order.UserID, order.Amount)
		resultChan <- order
	}
	close(resultChan)
}

func orderResult(resultChan <-chan Order, done chan bool) {
	for result := range resultChan {
		fmt.Printf("订单结果: ID=%d, 用户ID=%s, 金额=%d, 状态=%s\n", result.ID, result.UserID, result.Amount, result.Status)
	}
	done <- true
}

func main() {
	// orderChan := make(chan Order)
	// resultChan := make(chan Order)
	// done := make(chan bool)
	// go orderProduct(orderChan, 3)
	// go processOrder(orderChan, resultChan)
	// go orderResult(resultChan, done)
	// <-done // 阻塞在这里，直到 done channel 收到数据

	// 多路复用
	// 1. 谁先来，谁就先执行
	// ch1 := make(chan string)
	// ch2 := make(chan string)

	// go func() {
	// 	time.Sleep(time.Second * 2)
	// 	ch1 <- "hello"
	// }()

	// go func() {
	// 	time.Sleep(time.Second * 2)
	// 	ch2 <- "world"
	// }()

	// for i := 0; i < 2; i++ {
	// 	select {
	// 	case value := <-ch1:
	// 		fmt.Println("ch1:", value)
	// 	case value := <-ch2:
	// 		fmt.Println("ch2:", value)
	// 	case <-time.After(time.Second * 3):
	// 		fmt.Println("无消息")
	// 		return
	// 	}
	// }

	// 2. 定时器
	// 同时监听两个 channel，哪个有数据就处理哪个
	ticker := time.NewTicker(time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				fmt.Println("定时器:", t.Format("2006-01-02 15:04:05"))
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()
	time.Sleep(time.Second * 5)
	done <- true
}
