package main

import (
	"fmt"
	"sync"
	"time"
)

// 生产者/消费者模式：
//   • 生产者：负责生产数据，将数据发送到 channel
//   • 消费者：负责消费数据，从 channel 读取数据处理
//   • 通过 channel 解耦生产和消费，支持速率不一致的场景
//
// 扇入/扇出模式：
//   • 扇出：一个 channel 分发给多个 goroutine 处理（一产多消）
//   • 扇入：多个 channel 合并到一个 channel（多产一消）
//
// Pipeline 模式：将复杂任务分解为多个处理阶段，每个阶段通过 channel 连接，形成处理流水线。

// 生产者/消费者模式：
type Order struct {
	ID       int
	Amount   int
	Status   string
	CreateAt time.Time
}

// 生产者
func orderProduct(orderChan chan<- Order, number int) {
	defer close(orderChan)
	for i := 0; i < number; i++ {
		order := Order{
			ID:       i,
			Amount:   i * 100,
			Status:   "pending",
			CreateAt: time.Now(),
		}
		orderChan <- order
	}
}

// 消费者
func consumeOrder(orderChan <-chan Order, wg *sync.WaitGroup) {
	defer wg.Done()
	for order := range orderChan {
		time.Sleep(time.Millisecond * 500)
		order.Status = "processed"
		fmt.Println("消费订单:", order.ID, order.Status)
	}
}

// ----------------------------
// 生产者/消费者模式 main 函数
// func main() {
// 	var wg sync.WaitGroup
// 	// 创建订单 channel，容量为 5
// 	orderChan := make(chan Order, 5)

// 	// 启动消费者
// 	wg.Add(2)
// 	go consumeOrder(orderChan, &wg)
// 	go consumeOrder(orderChan, &wg)

// 	// 启动生产者
// 	go orderProduct(orderChan, 11)
// 	wg.Wait()
// }

// ----------------------------
// 扇出模式：一个生产者 -> 多个消费者（一产多消）
// 同一个 orderChan 分发给多个 worker 并发处理

func fanOut(orderChan <-chan Order, workerNum int) {
	var wg sync.WaitGroup
	for workerID := 0; workerID < workerNum; workerID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for order := range orderChan {
				time.Sleep(time.Millisecond * 100)
				fmt.Printf("[扇出] worker-%d 处理订单 ID:%d Amount:%d\n", workerID, order.ID, order.Amount)
			}
		}()
	}
	wg.Wait()
}

// ----------------------------
// 扇入模式：多个生产者 -> 一个消费者（多产一消）
// 多个 channel 合并到同一个 mergedChan

func merge(channels ...<-chan Order) <-chan Order {
	mergedChan := make(chan Order)
	var wg sync.WaitGroup
	for _, ch := range channels {
		wg.Add(1)
		c := ch
		go func() {
			defer wg.Done()
			for order := range c {
				mergedChan <- order
			}
		}()
	}
	go func() {
		wg.Wait()
		close(mergedChan)
	}()
	return mergedChan
}

func produceOrders(name string, start, count int) <-chan Order {
	ch := make(chan Order)
	go func() {
		defer close(ch)
		for i := start; i < start+count; i++ {
			ch <- Order{ID: i, Amount: i * 100, Status: "pending", CreateAt: time.Now()}
			fmt.Printf("[扇入] %s 生产订单 ID:%d\n", name, i)
			time.Sleep(time.Millisecond * 80)
		}
	}()
	return ch
}

// ----------------------------
// Pipeline 模式：订单处理流水线
// 阶段1: 生成订单 -> 阶段2: 校验订单 -> 阶段3: 计算折扣 -> 阶段4: 完成订单

// 阶段1：生成订单
func stageGenerate(count int) <-chan Order {
	out := make(chan Order)
	go func() {
		defer close(out)
		for i := 1; i <= count; i++ {
			out <- Order{ID: i, Amount: i * 100, Status: "pending", CreateAt: time.Now()}
			fmt.Printf("[阶段1-生成] 订单 ID:%d Amount:%d\n", i, i*100)
		}
	}()
	return out
}

// 阶段2：校验订单（过滤 Amount <= 0 的无效订单）
func stageValidate(in <-chan Order) <-chan Order {
	out := make(chan Order)
	go func() {
		defer close(out)
		for order := range in {
			if order.Amount <= 0 {
				fmt.Printf("[阶段2-校验] 订单 ID:%d 无效，跳过\n", order.ID)
				continue
			}
			order.Status = "validated"
			fmt.Printf("[阶段2-校验] 订单 ID:%d 校验通过\n", order.ID)
			out <- order
		}
	}()
	return out
}

// 阶段3：计算折扣（Amount > 300 享受 8 折）
func stageDiscount(in <-chan Order) <-chan Order {
	out := make(chan Order)
	go func() {
		defer close(out)
		for order := range in {
			if order.Amount > 300 {
				order.Amount = int(float64(order.Amount) * 0.8)
				fmt.Printf("[阶段3-折扣] 订单 ID:%d 享受8折，折后金额:%d\n", order.ID, order.Amount)
			} else {
				fmt.Printf("[阶段3-折扣] 订单 ID:%d 无折扣，金额:%d\n", order.ID, order.Amount)
			}
			out <- order
		}
	}()
	return out
}

// 阶段4：完成订单
func stageComplete(in <-chan Order) <-chan Order {
	out := make(chan Order)
	go func() {
		defer close(out)
		for order := range in {
			order.Status = "completed"
			fmt.Printf("[阶段4-完成] 订单 ID:%d 最终金额:%d 状态:%s\n", order.ID, order.Amount, order.Status)
			out <- order
		}
	}()
	return out
}

func main() {
	fmt.Println("======= 扇出示例 =======")
	fanOutChan := make(chan Order, 5)
	go func() {
		defer close(fanOutChan)
		for i := 0; i < 6; i++ {
			fanOutChan <- Order{ID: i, Amount: i * 100, Status: "pending", CreateAt: time.Now()}
		}
	}()
	fanOut(fanOutChan, 3) // 3 个 worker 并发消费

	fmt.Println("\n======= 扇入示例 =======")
	ch1 := produceOrders("生产者A", 0, 3)
	ch2 := produceOrders("生产者B", 10, 3)
	ch3 := produceOrders("生产者C", 20, 3)

	mergedChan := merge(ch1, ch2, ch3)
	for order := range mergedChan {
		fmt.Printf("[扇入] 消费者 收到订单 ID:%d Amount:%d\n", order.ID, order.Amount)
	}

	fmt.Println("\n======= Pipeline 示例 =======")
	// 将各阶段串联成流水线
	p1 := stageGenerate(5)
	p2 := stageValidate(p1)
	p3 := stageDiscount(p2)
	p4 := stageComplete(p3)

	// 收集最终结果
	for order := range p4 {
		_ = order
	}
}
