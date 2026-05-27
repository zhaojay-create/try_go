package main

import (
	"1-7-project-1/config"
	"1-7-project-1/other"
	"1-7-project-1/user"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var configFile string
	var port int
	// 这是命令行参数解析，从终端启动时传入
	flag.StringVar(&configFile, "config", "config.json", "path to config file")
	flag.IntVar(&port, "port", 8080, "port to listen")
	flag.Parse()
	config, err := config.LoadConfig(configFile)
	if err != nil {
		panic(err)
	}

	handler := user.NewUserHandler(user.NewUserStore())

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateUser(w, r)
		case http.MethodGet:
			handler.GetAllUser(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetUser(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	addr := fmt.Sprintf(":%d", config.Port)
	log.Printf("服务启动，监听 %s，日志级别: %s", addr, config.LogLevel)

	// other.RenderTemplate()
	// other.DemoCompression()
	other.DemoBufferIO()

	http.ListenAndServe(addr, nil)

}
