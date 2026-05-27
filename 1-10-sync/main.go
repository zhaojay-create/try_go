package main

import "fmt"

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

func main() {
	fmt.Println("Hello World")
}
