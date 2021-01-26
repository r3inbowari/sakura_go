package Sakura

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/*
 * 计划任务
 */
var LatestHome *goquery.Document // 最新主页快照
var cacheStart time.Time

func Init() {
	Blue("                                  ______             _            ")
	Blue("                                 |  ____|           (_)           ")
	Blue("  _ __ ___   ___  _ __ ___   ___ | |__   _ __   __ _ _ _ __   ___ ")
	Blue(" | '_ ` _ \\ / _ \\| '_ ` _ \\ / _ \\|  __| | '_ \\ / _` | | '_ \\ / _ \\")
	Blue(" | | | | | | (_) | | | | | | (_) | |____| | | | (_| | | | | |  __/")
	Blue(" |_| |_| |_|\\___/|_| |_| |_|\\___/|______|_| |_|\\__, |_|_| |_|\\___|")
	Blue("                                                __/ |             ")
	Blue("                                               |___/              ")

	Blue(" -- --- -- --- . -. --. .. -. . version: 1.0.0 author: r3inbowari")
	Blue("")
	Info("[SYS] system init")

	Info("[NETWORK TEST] " + RedirectURL("https://gss3.baidu.com/6LZ0ej3k1Qd3ote6lo7D0j9wehsv/tieba-smallvideo/607272_11d5cad2110530c892f7248946ebe51b.mp4"))
	Info("[NETWORK TEST] " + RedirectURL("http://quan.qq.com/video/1098_45b8f3ce393c72e8b8ebabee02fed632"))

	InitConfig()

	Info("权限申请: invenleey.oicp.net 验证通过")
	Info("日志等级: INFO")

	// server := NewServer()
	server := mux.NewRouter()

	server.HandleFunc("/hello", Hello)
	//server.Map("/rank", rank)
	//server.Map("/week", week)
	//server.Map("/last", last)
	//server.Map("/search", search)
	//server.Map("/detail", detail)
	//server.Map("/play", play)
	//server.useGlobalCORS()
	//
	//server.start()

	Info("解析服务启动: " + GetConfig().APIAddr)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Println("[INFO] 监听端口失败, 五秒后退出")
		time.Sleep(time.Second * 5)
		return
	}
	UseCache(true)
	HomepageSnapshotCron()



}
func Hello(w http.ResponseWriter, r *http.Request) {
	Warn("[3GPP] ")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Info("[ROUTE] " + r.RemoteAddr + " " + r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

/*
 * 重定向检查
 */
func CheckURL(baseHost string) string {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for {
		res, err := client.Get(baseHost)
		if err != nil {
			return baseHost
		}
		if res.StatusCode != 302 {
			return baseHost
		}
		baseHost = res.Header.Get("Location")
	}
}

func UseCache(sw bool) {
	if sw {
		Info("缓存服务: OPEN")
		c := cron.New()
		_ = c.AddFunc("0 */10 * * * ?", HomepageSnapshotCron)
		c.Start()
	}
}

func HomepageSnapshotCron() {
	cacheStart = time.Now()
	LatestHome, _ = goquery.NewDocument("http://www.yhdm.tv")
	spend := time.Now().UnixNano() - cacheStart.UnixNano()
	Info("缓存主页快照 耗时: " + strconv.Itoa(int(spend)/1e6) + "ms")
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

func play(w http.ResponseWriter, r *http.Request) {
	query := logRouter(r)
	bid, ok := query["bid"]
	if !ok || len(bid[0]) < 1 {
		log.Println("missing param 'bid'")
		requestFailed(w)
		return
	}

	num, ok := query["num"]
	if !ok || len(num[0]) < 1 {
		log.Println("missing param 'num'")
		requestFailed(w)
		return
	}

	p, err := strconv.Atoi(num[0])
	if err != nil {
		requestError(w)
		return
	}

	ab := getPlaySrc(bid[0], p)
	resultSucceed(w, CheckURL(ab), 1)
}

func detail(w http.ResponseWriter, r *http.Request) {
	query := logRouter(r)
	bid, ok := query["bid"]
	if !ok || len(bid[0]) < 1 {
		log.Println("missing param 'bid'")
		requestFailed(w)
		return
	}

	ab := detailBangumi(bid[0])
	resultSucceed(w, ab, 1)
}

func search(w http.ResponseWriter, r *http.Request) {
	query := logRouter(r)
	keyword, ok := query["keyword"]
	if !ok || len(keyword[0]) < 1 {
		log.Println("missing param 'keyword'")
		requestFailed(w)
		return
	}

	page, ok := query["page"]
	if !ok || len(page[0]) < 1 {
		ab, ba := searchBangumi(keyword[0], 1)
		resultSucceed(w, ab, ba)
	} else {
		p, err := strconv.Atoi(page[0])
		if err != nil {
			requestError(w)
			return
		}

		ab, ba := searchBangumi(keyword[0], p)
		resultSucceed(w, ab, ba)
	}
}

type RequestResult struct {
	Total   int         `json:"total"`
	Result  interface{} `json:"result"`
	Code    int         `json:"code"`
	Message string      `json:"msg"`
}

func logRouter(r *http.Request) url.Values {
	log.Println("[INFO] access route:", r.URL.Path)
	return r.URL.Query()
}

func rank(w http.ResponseWriter, r *http.Request) {
	logRouter(r)
	ab := rankUpdate()
	resultSucceed(w, ab, 1)
}

func last(w http.ResponseWriter, r *http.Request) {
	logRouter(r)
	ab := lastUpdate()
	resultSucceed(w, ab, 1)
}

func week(w http.ResponseWriter, r *http.Request) {
	query := logRouter(r)
	cal, ok := query["cal"]
	if !ok || len(cal[0]) < 1 {
		log.Println("missing param 'cal'")
		requestFailed(w)
		return
	}
	c, err := strconv.Atoi(cal[0])
	if c < 7 && err == nil {
		ab := weekUpdate(c)
		resultSucceed(w, ab, 1)
	} else {
		log.Println("max week list")
		requestError(w)
	}
	return
}

func requestFailed(w http.ResponseWriter) {
	var rq RequestResult
	rq.Result = nil
	rq.Total = 0
	rq.Code = 1
	rq.Message = "missing param"

	jsonStr, err := json.Marshal(rq)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	fmt.Fprintf(w, string(jsonStr))
}

func requestError(w http.ResponseWriter) {
	var rq RequestResult
	rq.Result = nil
	rq.Total = 0
	rq.Code = 1
	rq.Message = "error request"

	jsonStr, err := json.Marshal(rq)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	fmt.Fprintf(w, string(jsonStr))
}

func resultSucceed(w http.ResponseWriter, bangumi interface{}, total int) {
	var rq RequestResult
	rq.Result = bangumi
	rq.Total = total
	rq.Code = 0
	rq.Message = "succeed"
	jsonStr, err := json.Marshal(rq)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	fmt.Fprintf(w, string(jsonStr))
}

func rankUpdate() []SakuraBangumi {
	doc, _ := goquery.NewDocument("http://www.yhdm.tv")
	update := doc.Find(".side").Find(".pics").Find("li")
	var sbs []SakuraBangumi
	update.Each(func(i int, selection *goquery.Selection) {
		var sb SakuraBangumi
		u := selection.Find("img")
		sb.Title = u.AttrOr("alt", "null")
		sb.Cover = u.AttrOr("src", "null")
		sb.DetailURL = selection.Find("a").AttrOr("href", "null")

		sb.ID = strings.TrimLeft(sb.DetailURL, "/show/")
		sb.ID = strings.TrimRight(sb.ID, ".html")

		sbs = append(sbs, sb)
	})
	return sbs
}

func weekUpdate(w int) []SakuraBangumi {
	doc, _ := goquery.NewDocument("http://www.yhdm.tv")
	update := doc.Find(".side").Find(".tlist").Find("ul")
	var sbs []SakuraBangumi
	update.Each(func(i int, selection *goquery.Selection) {
		if w == i {
			selection.Find("li").Each(func(j int, sel *goquery.Selection) {
				var sb SakuraBangumi
				sb.DetailURL = sel.Find("a").Last().AttrOr("href", "null")
				sb.Title = sel.Find("a").Last().AttrOr("title", "null")
				sb.Update = sel.Find("span").Text()

				sb.ID = strings.TrimLeft(sb.DetailURL, "/show/")
				sb.ID = strings.TrimRight(sb.ID, ".html")

				sbs = append(sbs, sb)
			})
		}
	})
	return sbs
}

func GetData(url string) *goquery.Document {
	client := &http.Client{}
	resp, err := client.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	root, _ := html.Parse(resp.Body)
	return goquery.NewDocumentFromNode(root)

}

func lastUpdate() []SakuraBangumi {
	// doc, _ = goquery.NewDocument("http://www.yhdm.tv")
	start := time.Now().UnixNano()
	update := LatestHome.Find(".firs").Find(".img").First().Find("li")
	var sbs []SakuraBangumi
	update.Each(func(i int, selection *goquery.Selection) {
		var sb SakuraBangumi
		sb.Title = selection.Find("img").AttrOr("alt", "null")
		sb.Cover = selection.Find("img").AttrOr("src", "null")
		sb.Update = selection.Find("p").Last().Text()
		sb.DetailURL = selection.Find("a").AttrOr("href", "null")
		sb.ID = strings.TrimLeft(sb.DetailURL, "/show/")
		sb.ID = strings.TrimRight(sb.ID, ".html")

		sbs = append(sbs, sb)
	})
	end := time.Now().UnixNano()
	log.Println("[INFO] 耗时: ", (end-start)/100000, "ms")
	return sbs
}

func getPlaySrc(id string, num int) string {
	doc := GetSearchResult("http://www.yhdm.tv/v/" + id + "-" + strconv.Itoa(num) + ".html")
	playSrc := doc.Find(".bofang").Find("div").AttrOr("data-vid", "null")

	defer func() {
		if e := recover(); e != nil {
			log.Println("[INFO] 地址解析失败: url")
		}
	}()

	inx := strings.LastIndex(playSrc, "$")
	return playSrc[:inx]
}

func detailBangumi(id string) SakuraBangumi {
	doc := GetSearchResult("http://www.yhdm.tv/show/" + id + ".html")
	// doc, _ := goquery.NewDocument("http://www.yhdm.tv/show/" + id + ".html")
	pInfo := doc.Find(".sinfo").Find("p")
	var sb SakuraBangumi
	if pInfo.Size() > 1 {
		sb.Alias = strings.TrimLeft(pInfo.First().Text(), "别名:")
	}
	sb.Detail = doc.Find(".info").Text()
	sb.Count = doc.Find(".movurl").Find("li").Size()

	sb.Title = doc.Find(".thumb").Find("a").AttrOr("title", "null")
	sb.Cover = doc.Find("img").AttrOr("src", "null")

	sb.DetailURL = "/show/" + id + ".html"
	sb.Update = pInfo.Last().Text()

	doc.Find(".sinfo").Find("span").First().Next().Next().Find("a").Each(func(i int, selection *goquery.Selection) {
		sb.Type = append(sb.Type, selection.AttrOr("href", "null"))
	})

	sb.ID = id
	return sb
}

func searchBangumi(keyword string, page int) ([]SakuraBangumi, int) {
	doc := GetSearchResult("http://www.yhdm.tv/search/" + keyword + "/?page=" + strconv.Itoa(page))
	a := doc.Find(".lpic").Find("li")
	var sbs []SakuraBangumi
	a.Each(func(i int, selection *goquery.Selection) {
		var sb SakuraBangumi
		sb.Title = selection.Find("a").Text()
		sb.Cover = selection.Find("img").AttrOr("src", "null")
		sb.Detail = selection.Find("p").Text()
		sb.DetailURL = selection.Find("a").AttrOr("href", "null")
		sb.Update = selection.Find("span").First().Text()

		selection.Find("span").Find("a").Each(func(i int, sel *goquery.Selection) {
			sb.Type = append(sb.Type, sel.AttrOr("href", "null"))
		})

		sb.ID = strings.TrimLeft(sb.DetailURL, "/show/")
		sb.ID = strings.TrimRight(sb.ID, ".html")
		sbs = append(sbs, sb)
	})
	strCount := doc.Find("#totalnum").Text()
	if strCount != "" {
		count, err := strconv.Atoi(strings.TrimRight(strCount, "条"))
		if err == nil {
			return sbs, (count + 19) / 20
		}
	}
	return sbs, 1
}
