package main

import (
	"bufio"
	"github.com/cute-angelia/go-utils/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go-deploy/config"
	"go-deploy/internal/controller"
	"log"
	"net"
	"net/http"
	"time"
)

const TYPE_SITE = "go-deploy-server"

func main() {
	// config
	config.InitConfig()

	// 日志
	logger.NewLogger(TYPE_SITE, !config.C.Debug)

	// router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.ThrottleBacklog(20, 500, time.Second))

	// 中间件
	r.Route("/", func(r chi.Router) {
		r.Mount("/", controller.Index{}.Routes())

		r.Mount("/list", controller.List{}.Routes())
		r.Mount("/deply", controller.Deply{}.Routes())
		r.Mount("/rollback", controller.Rollback{}.Routes())
		r.Mount("/showlog", controller.ShowLog{}.Routes())
	})

	// Ping
	for addr := range config.C.UniqAddr {
		go Ping(addr)
	}

	log.Println(TYPE_SITE + "启动成功~ " + config.C.ListenHttp)
	if err := http.ListenAndServe(config.C.ListenHttp, r); err != nil {
		log.Println("启动错误：", err)
	}
}

func Ping(addr string) {
	for {
		func() {
			//connect to this socket
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				log.Println("Error connect to client:", err)
				return
			}

			//connect success
			setClientOnlineStatus(addr, true)

			//remote client closed
			defer func() {
				setClientOnlineStatus(addr, false)
				conn.Close()
			}()

			//read message from client
			go func(conn net.Conn) {
				defer conn.Close()
				for {
					message, err := bufio.NewReader(conn).ReadString('\n')
					if err != nil {
						log.Println("Client closed", err.Error())
						return
					}
					if config.C.Debug {
						log.Print(conn.RemoteAddr(), " -> Message Received from client:", message)
					}
				}
			}(conn)

			ticker := time.Tick(time.Second * 15)
			for {
				select {
				case <-ticker:
					conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
					_, err := conn.Write([]byte("PING\n"))
					if err != nil {
						log.Println("Error writing to stream:", err)
						return
					}
				}
			}
		}()
		time.Sleep(time.Second * 5)
	}
}

//set client online or offline
func setClientOnlineStatus(addr string, online bool) {
	for key, app := range config.C.Apps {
		for k, node := range app.Node {
			if node.Addr == addr {
				config.C.Apps[key].Node[k].Online = online
			}
		}
	}
}
