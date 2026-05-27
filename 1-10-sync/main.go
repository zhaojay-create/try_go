package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// sync.mutex 互斥锁
// 同一时间，只有一个 goroutine 能够访问资源
// Lock 和 UnLock
// 使用后必须释放，否则会造成死锁

// sync.RWMutex 读写锁
// 允许多个读，一个写
// 读锁：RLock 和 RUnlock
// 写锁：Lock 和 Unlock
// 使用后必须释放，否则会造成死锁
// 写锁了，不能读
// 读锁了，不能写，但是能读

// sync.WaitGroup 等待组
// 用于等待一组 goroutine 执行完成
// Add 增加等待的 goroutine 数量
// Done 减少等待的 goroutine 数量
// Wait 等待所有 goroutine 执行完成

type Inventory struct {
	stock   int
	rwMutex sync.RWMutex
}

func (v *Inventory) getStock() int {
	v.rwMutex.RLock()
	defer v.rwMutex.RUnlock()

	return v.stock
}

func (v *Inventory) deductStock(quantity int) bool {
	v.rwMutex.Lock()
	defer v.rwMutex.Unlock()
	if v.stock < quantity {
		fmt.Printf("库存不足")
		return false
	}
	time.Sleep(time.Millisecond * 100)
	v.stock -= quantity
	return true
}

type Counter struct {
	value int64
}

func (c *Counter) Inc() {
	atomic.AddInt64(&c.value, 1)
}

func (c *Counter) Dec() {
	atomic.AddInt64(&c.value, -1)
}

func (c *Counter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}

func main() {
	// inventory := &Inventory{stock: 100}

	var wg sync.WaitGroup

	// for i := 0; i < 50; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		inventory.deductStock(1)
	// 	}()
	// }

	// for i := 0; i < 50; i++ {
	// 	time.Sleep(time.Millisecond * 100)
	// 	fmt.Printf("剩余数量: %d \n", inventory.getStock())
	// }
	// wg.Wait()

	// fmt.Printf("所有协程执行完毕，剩余数量: %d\n", inventory.getStock())

	var counter Counter
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Inc()
		}()
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Dec()
		}()
	}

	wg.Wait()
	fmt.Printf("所有协程执行完毕，点赞数量: %d\n", counter.Get())
}
