package Sakura

import (
	"github.com/gorilla/mux"
	"github.com/wuwenbao/gcors"
	"net/http"
	"time"
)

type Server struct {
	router *mux.Router
	cors   bool
}

func NewServer() Server {
	r := mux.NewRouter()
	return Server{router: r}
}

func (s *Server) Map(path string, f func(http.ResponseWriter,
	*http.Request)) *Server {
	Info("[RT] add route path -> " + path)
	s.router.HandleFunc(path, f)
	return s
}

func (s *Server) start() {
	Info("[SERVER] Listen on: " + GetConfig().APIAddr)

	if s.cors {
		cors := gcors.New(
			s.router,
			gcors.WithOrigin("*"),
			gcors.WithMethods("POST, GET, PUT, DELETE, OPTIONS"),
			gcors.WithHeaders("Authorization"),
		)

		err := http.ListenAndServe(GetConfig().APIAddr, cors)
		if err != nil {
			Fail("[INFO] Listen failed, exit after 5 second")
			time.Sleep(time.Second * 5)
			return
		}
	} else {
		err := http.ListenAndServe(GetConfig().APIAddr, s.router)
		if err != nil {
			Fail("[INFO] Listen failed, exit after 5 second")
			time.Sleep(time.Second * 5)
			return
		}
	}
}

func (s *Server) useGlobalCORS() {
	Info("[SERVER] Global CORS OPEN")
	s.cors = true

}

func (s *Server) useLog() {
	Info("[SERVER] Router log OPEN")
	s.router.Use(loggingMiddleware)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Info("[ROUTE] " + r.RemoteAddr + " " + r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
