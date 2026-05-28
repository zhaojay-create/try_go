package main

import "github.com/gin-gonic/gin"

func main() {
	// 1. 初始化 gin 引擎
	r := gin.Default()
	// 2. 定义路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 3. 启动服务
	r.Run()
}
