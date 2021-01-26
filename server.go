package Sakura

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Server struct {
	router *mux.Router
}

func NewServer() Server {
	r := mux.NewRouter()
	// domain(r)

	return Server{router: r}
}

func (s *Server) Map(path string, f func(http.ResponseWriter,
	*http.Request)) {
	Info("[RT] add route path -> " + path)
	s.router.HandleFunc(path, f)
}

func (s *Server) start() {
	Info("解析服务启动: " + GetConfig().APIAddr)
	err := http.ListenAndServe(GetConfig().APIAddr, nil)
	if err != nil {
		log.Println("[INFO] 监听端口失败, 五秒后退出")
		time.Sleep(time.Second * 5)
		return
	}
}

func (s *Server) useGlobalCORS() {
	Info("[CORS] OPEN")
	s.router.Use(loggingMiddleware)
}
