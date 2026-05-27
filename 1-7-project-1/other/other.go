package other

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"strings"
	"text/template"
)

// Go 的访问控制规则：

// 大写开头 → exported，包外可以访问
// 小写开头 → unexported，只能包内使用

func RenderTemplate() {
	const templateText = `
	用户信息:
	姓名: {{.Name}}
	邮箱: {{.Email}}
	创建时间: {{.CreateAt}}
	`

	type User struct {
		Name     string
		Email    string
		CreateAt string
	}

	user := User{
		Name:     "张三",
		Email:    "zhangsan@example.com",
		CreateAt: "2025-10-13",
	}

	temp, err := template.New("user").Parse(templateText)
	if err != nil {
		panic(err)
	}
	err = temp.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}

func DemoCompression() {
	original := "这是一段需要压缩的文本," + strings.Repeat("hello world", 100)

	var compressed bytes.Buffer

	gzWriter := gzip.NewWriter(&compressed)

	_, err := gzWriter.Write([]byte(original))

	if err != nil {
		panic(err)
	}

	gzWriter.Close()

	println("原始数据大小:", len(original))
	println("压缩后数据大小:", compressed.Len())

	gzReader, err := gzip.NewReader(&compressed)
	if err != nil {
		panic(err)
	}
	defer gzReader.Close()

	decompressed, err := io.ReadAll(gzReader)
	if err != nil {
		panic(err)
	}

	println("解压后数据大小:", len(decompressed))
	println("解压后数据:", string(decompressed))
}

func DemoBufferIO() {
	filePath := "demo.txt"

	// 创建并写入文件
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString("你好，这是写入的内容\n第二行内容\n")
	if err != nil {
		panic(err)
	}
	file.Close()

	// 读取文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	println("读取到的内容:")
	println(string(content))
}
