package main

import "fmt"

// demo1: 泛型函数 - 打印任意类型
// T 是类型参数，any 表示 T 可以是任何类型
func Print[T any](val T) {
	fmt.Println(val)
}

// demo2: 泛型函数 - 返回两个值中较大的一个
// comparable 约束 T 必须支持 > 运算符（int/float/string）
// Number 定义的是允许的类型范围：
type Number interface {
	~int | ~float64 | ~string
}

func Max[T Number](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// demo3: 泛型函数 - 过滤切片
func Filter[T any](slice []T, fn func(T) bool) []T {
	result := []T{}
	for _, v := range slice {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// demo4: 泛型结构体 - 通用栈
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

func main() {
	// demo1: 同一个函数，传不同类型
	Print[string]("hello")
	Print[int](42)
	Print(3.14) // 类型可以自动推断，不用显式写 [float64]

	// demo2: 比较大小
	fmt.Println(Max(3, 5))              // 5
	fmt.Println(Max(3.14, 2.71))        // 3.14
	fmt.Println(Max("apple", "banana")) // banana

	// demo3: 过滤切片
	nums := []int{1, 2, 3, 4, 5, 6}
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println(evens) // [2 4 6]

	words := []string{"go", "python", "java", "rust"}
	short := Filter(words, func(s string) bool { return len(s) <= 3 })
	fmt.Println(short) // [go]

	// demo4: 泛型栈
	intStack := Stack[int]{}
	intStack.Push(1)
	intStack.Push(2)
	intStack.Push(3)
	val, _ := intStack.Pop()
	fmt.Println(val) // 3

	strStack := Stack[string]{}
	strStack.Push("hello")
	strStack.Push("world")
	s, _ := strStack.Pop()
	fmt.Println(s) // world
}
