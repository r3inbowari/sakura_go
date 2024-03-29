package Sakura

import (
	"context"
	"fmt"
	"github.com/aliyun/fc-runtime-go-sdk/fc"
	"github.com/aliyun/fc-runtime-go-sdk/fccontext"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/wuwenbao/gcors"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var BiliServer *Server

type Server struct {
	router *mux.Router
	s      *http.Server
}

func ServerMap() {
	Log.Info("[BCS] sakura CLI is running")
	BiliServer = NewServer()
	BiliServer.Map("/version", HandleVersion)
	BiliServer.Map("/play/{id}/{sid}/{nid}", HandleGetPlayer)
	BiliServer.Map("/rank", HandleRank)
	BiliServer.Map("/new", HandleNew)
	BiliServer.Map("/weeks/{index}", HandleWeeks)
	BiliServer.Map("/search", HandleSearch)
	BiliServer.Map("/detail/{id}", HandleDetail)
	BiliServer.router.Use(loggingMiddleware)
}

func CLIApplication() {
	ServerMap()
	if Up.BuildMode == "AliyunFC" {
		fc.StartHttp(AFCServer)
	} else {
		NormalServer()
	}
}

func NormalServer() {
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

// AFCServer 阿里云函数计算服务
func AFCServer(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	lc, _ := fccontext.FromContext(ctx)
	fmt.Printf("上下文: %#v\n", lc)
	fmt.Printf("请求头: %#v\n", req.Header)
	fmt.Printf("路由点: %#v\n", req.URL.String())
	fmt.Printf("方法名: %#v\n", req.Method)
	BiliServer.router.ServeHTTP(w, req)
	return nil
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
		start := time.Now().UnixNano()
		if Up.BuildMode == "AliyunFC" {
			DMSakuraHomepageSnapshotCron()
		}
		unescape, err := url.QueryUnescape(r.RequestURI)
		if err != nil {
			Log.Error("error route")
			return
		}
		next.ServeHTTP(w, r)
		end := time.Now().UnixNano()
		Log.WithFields(logrus.Fields{"time": fmt.Sprintf("%dms", (end-start)/100000)}).Infof("[BSC] route %s -> %s", r.RemoteAddr, unescape)
	})
}

func HandleVersion(w http.ResponseWriter, r *http.Request) {
	ResponseCommon(w, Up.VersionStr+" "+Up.ReleaseTag, "ok", 1, http.StatusOK, 0)
}
