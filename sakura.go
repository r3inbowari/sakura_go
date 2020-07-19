package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron"
	"strconv"
	"time"
)

/*
 * 计划任务
 */
var LatestHome *goquery.Document // 最新主页快照
var cacheStart time.Time

func Init() {
	INFO(RedirectURL("https://gss3.baidu.com/6LZ0ej3k1Qd3ote6lo7D0j9wehsv/tieba-smallvideo/607272_11d5cad2110530c892f7248946ebe51b.mp4"))
	INFO(RedirectURL("http://quan.qq.com/video/1098_45b8f3ce393c72e8b8ebabee02fed632"))
	INFO("权限申请: invenleey.oicp.net 验证通过")
	INFO("日志等级: INFO")
	INFO("缓存服务: OPEN")
	HomepageCron()

	c := cron.New()
	_ = c.AddFunc("0 */10 * * * ?", HomepageCron)
	c.Start()
}

func HomepageCron() {
	cacheStart = time.Now()
	LatestHome, _ = goquery.NewDocument("http://www.yhdm.tv")
	spend := time.Now().UnixNano() - cacheStart.UnixNano()
	INFO("缓存主页快照 耗时: " + strconv.Itoa(int(spend)/1e6) + "ms")
}

type SakuraBangumi struct {
	Title     string   `json:"title"`
	DetailURL string   `json:"detailURL"`
	Cover     string   `json:"cover"`
	Detail    string   `json:"detail"`
	Update    string   `json:"update"`
	Count     int      `json:"count"`
	Type      []string `json:"type"`
	Alias     string   `json:"alias"`
	ID        string   `json:"id"`
}
