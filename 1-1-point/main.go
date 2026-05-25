package main

import "fmt"

type User struct {
	Id   int
	Name string
	Age  int
}

// 不需要指针：读取操作（不需要修改原变量）
func (u User) GetInfo() string {
	return fmt.Sprintf("%s (ID: %d, Age: %d)", u.Name, u.Id, u.Age)
}

// 需要指针：修改操作（需要改变原变量的值）
func (u *User) SetAge(age int) {
	u.Age = age
}

func main() {
	user := User{
		Id:   1,
		Name: "张三",
		Age:  18,
	}

	// 1. 不需要指针的方法：读取
	info := user.GetInfo()
	fmt.Println("原始:", info)

	// 2. 需要指针的方法：修改
	user.SetAge(25)
	fmt.Println("修改后:", user.GetInfo())
}
