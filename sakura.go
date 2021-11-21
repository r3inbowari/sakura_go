package Sakura

import (
	"github.com/PuerkitoBio/goquery"
	"time"
)

// LatestHome /*
var LatestHome *goquery.Document // 最新主页快照
var cacheStart time.Time
var RDB *Snapshot

func Init() {
	Log.Info("[SYS] system init")
	Log.Info("[NETWORK TEST] " + RedirectURL("https://gss3.baidu.com/6LZ0ej3k1Qd3ote6lo7D0j9wehsv/tieba-smallvideo/607272_11d5cad2110530c892f7248946ebe51b.mp4"))
	Log.Info("[NETWORK TEST] " + RedirectURL("http://quan.qq.com/video/1098_45b8f3ce393c72e8b8ebabee02fed632"))

	RDB = InitCacheService()
	RDB.UseCache()

	CLIApplication()
}

type RequestResult struct {
	Total int         `json:"total"`
	Data  interface{} `json:"data"`
	Code  int         `json:"code"`
	Message string      `json:"msg"`
}
