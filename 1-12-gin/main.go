package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化 gin 引擎
	r := gin.Default()
	// 全局中间件，对所有路由生效
	r.Use(func(c *gin.Context) {
		fmt.Println("global middleware")
		c.Next()
	})
	// 2. 定义路由
	// 3.1 JSON 响应
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 3.2 响应 html
	// 通配符匹配批量加载
	// r.LoadHTMLGlob("html/*")

	// 指定具体文件路径
	r.LoadHTMLFiles("html/index.html")
	r.GET("/html", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// 3.3 响应文件
	r.GET("/file", func(c *gin.Context) {
		// 表示文件响应头，唤起浏览器的下载
		c.Header("Content-Type", "application/octet-stream")
		// 设置下载下来的文件名字
		c.Header("Content-Disposition", "attachment; filename*=UTF-8''go%E6%B5%8B%E8%AF%95")
		c.File("html/file.go")
	})
	// 3.3.1 下载方式
	// 后端直接返回二进制流
	// 前端 → GET /download → 后端直接返回文件内容(blob) → 前端构造a标签 → 下载
	// 	文件内容先加载进浏览器内存再下载，大文件会占用大量内存
	// 需要等整个文件加载完才能开始下载
	// 无法支持断点续传
	// 3.3.2 下载方式 后端返回临时URL
	// 前端 → POST /get-download-url → 后端返回临时URL（如 https://cdn.xxx.com/file?token=xxx&expires=1234）
	// 前端 → 用这个URL构造a标签 → 浏览器直接向CDN/对象存储请求文件 → 下载

	// 3.4 静态文件
	r.Static("stc", "static")
	r.StaticFile("abc", "static/abc.txt")

	// 4.1 查询参数的请求
	// http://localhost:8080/goods?name=123&age=18&key=abc&key=qwe
	// 	{
	//	 "age": "18",
	//	 "code": 200,
	//	 "keys": [
	//	 "[abc]",
	//	 "qwe"
	//	 ],
	//	 "name": "123",
	//	 "status": true
	// }
	r.GET("/goods", func(c *gin.Context) {
		// 获取查询参数
		name := c.Query("name")
		age := c.DefaultQuery("age", "18")
		keys := c.QueryArray("key")
		c.JSON(200, gin.H{
			"code":   200,
			"status": true,
			"name":   name,
			"age":    age,
			"keys":   keys,
		})
	})

	// 4.2 动态参数的请求
	// http://localhost:8080/user/123/order/456
	r.GET("/user/:id/order/:orderId", func(c *gin.Context) {
		id := c.Param("id")
		orderId := c.Param("orderId")
		c.JSON(200, gin.H{
			"code":    200,
			"id":      id,
			"orderId": orderId,
		})
	})
	// 4.3 表单参数的请求
	// POST http://localhost:8080/login  body: username=admin&password=123
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.DefaultPostForm("password", "000000")
		age, ok := c.GetPostForm("age")
		if !ok {
			age = "0"
		}
		c.JSON(200, gin.H{
			"code":     200,
			"username": username,
			"password": password,
			"age":      age,
		})
	})

	// 5.1 查询参数的绑定
	type SearchQuery struct {
		Name string `form:"name" binding:"required"`
		Age  int    `form:"age"  binding:"min=1,max=120"`
	}
	// GET http://localhost:8080/search?name=123&age=18
	r.GET("/search", func(c *gin.Context) {
		var query SearchQuery
		if err := c.ShouldBindQuery(&query); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"query": query})
	})

	// 5.2 路径参数的绑定
	type UserOrderUri struct {
		ID      int `uri:"id"      binding:"required,min=1"`
		OrderID int `uri:"orderId" binding:"required,min=1"`
	}
	// GET http://localhost:8080/u/123/order/456
	r.GET("/u/:id/order/:orderId", func(c *gin.Context) {
		var uri UserOrderUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"uri": uri})
	})

	// 5.3 JSON 参数的绑定
	type LoginBody struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}
	// POST http://localhost:8080/loginJson  body: {"username":"admin","password":"123456"}
	r.POST("/loginJson", func(c *gin.Context) {
		var body LoginBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"body": body})
	})

	// 5.4 表单参数的绑定
	type RegisterForm struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required,min=6"`
		Age      int    `form:"age"      binding:"min=1,max=120"`
	}
	// POST http://localhost:8080/register  body: username=admin&password=123456&age=18
	r.POST("/register", func(c *gin.Context) {
		var form RegisterForm
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"form": form})
	})

	// 6. gin 路由 group
	v1 := r.Group("/v1")
	// 路由分组可以统一的加入中间间
	v1.Use(func(c *gin.Context) {
		fmt.Println("v1 group middleware")
		c.Next()
	})
	v1Group(v1)

	// 7.1 局部中间件
	M1 := func(c *gin.Context) {
		fmt.Println("m1 middleware start")
		c.Next()
		fmt.Println("m1 middleware end")
	}
	M2 := func(c *gin.Context) {
		fmt.Println("m2 middleware start")
		c.Next()
		fmt.Println("m2 middleware end")
	}
	r.GET("/middleware/ping", M1, M2, func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "v1pong"})
	})

	// 7.2 全局中间件，是对路由组的
	// 定义路由组
	auth := r.Group("/auth")
	auth.Use(func(c *gin.Context) {
		fmt.Println("auth group middleware")
		c.Next()
	})
	auth.GET("/login", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "auth login"})
	})
	auth.GET("/logout", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "auth logout"})
	})

	r.Run()

}

func v1Group(r *gin.RouterGroup) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "v1pong"})
	})
}
