package main

import "fmt"

// 数组：数组是固定长度，值传递，当函数参数是数组时，会创建一个副本
// slice：切片是动态长度，引用传递，当函数参数是切片时，不会创建副本,
// 传递时复制切片头（含指针），修改元素会影响原切片，但 append 扩容后不影响原切片。

// 切片底层是一个结构体，包含3个字段：

// ptr（指针） — 指向底层数组第一个元素的内存地址
// len（长度） — 当前切片中有效元素的个数
// cap（容量） — 底层数组从该指针位置起最多能存多少元素

type CartItem struct {
	Id      int     // 商品ID
	Name    string  // 商品名称
	Price   float64 // 价格
	Quanity int     // 数量
}

func main() {
	// 定义切片
	var cart []CartItem
	cart = append(cart, CartItem{Id: 1, Name: "1号商品", Price: 30.0, Quanity: 1})
	cart = append(cart, CartItem{Id: 2, Name: "2号商品", Price: 20.0, Quanity: 2})

	total := 0.0
	totalQuanity := 0
	for _, item := range cart {
		total += item.Price * float64(item.Quanity)
		totalQuanity += item.Quanity
		fmt.Printf("商品id: %d, 名称: %s, 价格: %.2f, 数量: %d\n", item.Id, item.Name, item.Price, item.Quanity)
	}
	fmt.Printf("总价: %.2f\n", total)
	fmt.Printf("商品数量: %d\n", totalQuanity)

	// 模拟商品移除
	// len(cart) : 返回切片 cart 中元素的数量
	if len(cart) > 0 {
		cart = cart[1:]
	}

	total = 0.0
	totalQuanity = 0
	for _, item := range cart {
		total += item.Price * float64(item.Quanity)
		totalQuanity += item.Quanity
		fmt.Printf("商品id: %d, 名称: %s, 价格: %.2f, 数量: %d\n", item.Id, item.Name, item.Price, item.Quanity)
	}
	fmt.Printf("移除商品后，总价: %.2f\n", total)
	fmt.Printf("移除商品后，商品数量: %d\n", totalQuanity)
}
