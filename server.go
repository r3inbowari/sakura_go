package Sakura

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/wuwenbao/gcors"
	"net/http"
	"os"
	"strings"
	"time"
)

var BiliServer *Server

type Server struct {
	router *mux.Router
	s      *http.Server
}

func CLIApplication() {
	Log.Info("[BCS] MEIWOBUXING CLI PACKAGER is running")
	BiliServer = NewServer()
	BiliServer.Map("/version", HandleVersion)

	BiliServer.router.Use(loggingMiddleware)
	err := BiliServer.start()
	if strings.HasSuffix(err.Error(), "normally permitted.") || strings.Index(err.Error(), "bind") != -1 {
		Log.WithFields(logrus.Fields{"err": err.Error()}).Error("[BCS] Only one usage of each socket address is normally permitted.")
		Log.Info("[BCS] EXIT 1002")
		os.Exit(1002)
	}
	// goroutine block here not need sleep
	Log.WithFields(logrus.Fields{"err": err.Error()}).Info("[BCS] Service will be terminated soon")
	time.Sleep(time.Second * 10)
}

func Shutdown(ctx context.Context) {
	BiliServer.Shutdown(ctx)
}

func NewServer() *Server {
	r := mux.NewRouter()
	Log.Info("[BSC] Global CORS OPEN")

	cors := gcors.New(
		r,
		gcors.WithOrigin("*"),
		gcors.WithMethods("POST, GET, PUT, DELETE, OPTIONS"),
		gcors.WithHeaders("Authorization"),
	)

	retServer := &http.Server{
		Addr:    GetConfig(false).APIAddr,
		Handler: cors,
	}
	return &Server{router: r, s: retServer}
}

func (s *Server) start() error {
	Log.Info("[BCS] Listened on " + GetConfig(false).APIAddr)
	return s.s.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) {
	if s.s != nil {
		Log.Info("[BSC] releasing server now...")
		err := s.s.Shutdown(ctx)
		if err != nil {
			Log.Error("[BSC] Shutdown failed")
			Log.Info("[BCS] EXIT 1002")
			os.Exit(1011)
		}
		Log.Info("[BSC] release completed")
	}
}

func (s *Server) Map(path string, f func(http.ResponseWriter,
	*http.Request), method ...string) *Server {
	if len(method) == 1 {
		Log.Info("[BSC] add route path [" + method[0] + "] -> " + path)
		s.router.HandleFunc(path, f).Methods(method[0])
	} else {
		Log.Info("[BSC] add route path [ALL] -> " + path)
		s.router.HandleFunc(path, f)
	}
	return s
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Log.Info("[BSC] route" + r.RemoteAddr + " " + r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func HandleVersion(w http.ResponseWriter, r *http.Request) {
	ResponseCommon(w, Up.VersionStr+" "+Up.ReleaseTag, "ok", 1, http.StatusOK, 0)
}