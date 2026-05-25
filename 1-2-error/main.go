package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

// defer 语句用于延迟执行，通常用于资源清理
// defer 语句会在函数返回前执行，即使函数中发生了 panic / return

// ==================== 1. error 类型 ====================
// error 是 Go 的内置接口：type error interface { Error() string }
// 用于显式返回错误，强制调用者处理

// 自定义错误
type BusinessError struct {
	Code    int
	Message string
	Time    time.Time
}

// 实现 go 的 error 接口
func (e *BusinessError) Error() string {
	return fmt.Sprintf("错误代码: %d, 错误信息: %s, 时间: %s", e.Code, e.Message, e.Time.Format("2006-01-02 15:04:05"))
}

func queryDataBase(userId int) (string, error) {
	if userId < 0 {
		return "", &BusinessError{
			Code:    1002,
			Message: "用户ID不能为负数",
			Time:    time.Now(),
		}
	}
	if userId == 10 {
		return "张三", nil
	}

	return "", &BusinessError{
		Code:    1001,
		Message: "用户不存在",
		Time:    time.Now(),
	}
}

// ==================== 2. defer 语句 ====================
// defer: 延迟执行，函数返回前执行（无论是否出错）
// 用途：资源清理（关闭文件、解锁、释放连接等）
// 特点：LIFO 顺序（后 defer 的先执行）

// readFile 演示 defer 关闭文件
// 无论函数正常返回还是出错，file.Close() 都会执行
func readFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	// defer 在 return 前执行，确保文件一定被关闭
	defer file.Close()

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	return string(content), nil
}

// ==================== 3. panic / recover ====================
// panic: 立即停止当前函数，开始向上层栈回溯
// recover: 捕获 panic，阻止程序崩溃，必须在 defer 中调用
// 用途：处理不可恢复的错误或防止程序崩溃

// safeAccess 演示捕获数组越界 panic
func safeAccess(arr []int, index int) (value int, err error) {
	// defer + recover 捕获 panic
	defer func() {
		if r := recover(); r != nil {
			// 捕获到 panic，转为正常错误返回
			err = fmt.Errorf("索引越界: %v", r)
		}
	}()

	// 如果 index 越界，这里会触发 panic
	return arr[index], nil
}

func main() {
	fmt.Println("===== 1. error 类型 =====")
	// Go 显式错误处理，每个可能出错的函数都返回 error
	name, err := queryDataBase(10)
	if err != nil {
		fmt.Println("查询失败:", err)
	} else {
		fmt.Println("查询成功:", name)
	}

	fmt.Println("\n===== 2. defer 语句（文件读取示例）=====")
	// 读取自身源代码文件（defer 确保关闭）
	content, err := readFile("1-2-error/main.go")
	if err != nil {
		fmt.Println("读取失败:", err)
	} else {
		// 只打印前 100 个字符
		preview := content
		if len(content) > 100 {
			preview = content[:100] + "..."
		}
		fmt.Println("文件内容预览:", preview)
	}

	// 尝试读取不存在的文件（defer 仍会执行清理）
	fmt.Println("\n尝试读取不存在的文件...")
	_, err = readFile("not_exist.txt")
	if err != nil {
		fmt.Println("预期中的错误:", err)
	}

	fmt.Println("\n===== 3. panic / recover（数组越界示例）=====")
	nums := []int{10, 20, 30}

	// 场景 A：正常访问
	val, err := safeAccess(nums, 1)
	if err != nil {
		fmt.Println("访问出错:", err)
	} else {
		fmt.Printf("nums[1] = %d\n", val)
	}

	// 场景 B：越界访问（触发 panic，但被 recover 捕获）
	val, err = safeAccess(nums, 10)
	if err != nil {
		fmt.Println("访问出错（panic 被捕获）:", err)
	} else {
		fmt.Printf("nums[10] = %d\n", val)
	}
	fmt.Println("程序继续执行，没有崩溃！")

	// 场景 C：不捕获，程序崩溃（注释掉 recover 就会看到）
	// _ = nums[10] // 直接越界，程序崩溃
}
