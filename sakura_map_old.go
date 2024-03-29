package Sakura

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strconv"
)

// 排行榜
func rank(w http.ResponseWriter, r *http.Request) {
	ru := RankUpdate()
	resultSucceed(w, ru, 1)
}

// 最新内容
func last(w http.ResponseWriter, r *http.Request) {
	ab := LastUpdate()
	resultSucceed(w, ab, 1)
}

// 周内容
func week(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dayIndex := vars["index"]
	c, err := strconv.Atoi(dayIndex)
	if c < 7 && err == nil {
		ab := WeekUpdate(c)
		resultSucceed(w, ab, 1)
	} else {
		Log.Warn("[API] apply error param")
		requestError(w)
	}
	return
}

// 关键字搜索
func search(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	keyword := r.FormValue("keyword")
	page := r.FormValue("page")

	p, err := strconv.Atoi(page)
	if err == nil {
		ab, ba := SearchBangumi(keyword, p+1)
		resultSucceed(w, ab, ba)
	} else {
		Log.Warn("[API] apply error param")
		requestError(w)
	}
}

// 详细页
func detail(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	bid := r.FormValue("bid")
	res := DetailBangumi(bid)
	resultSucceed(w, res, 1)
}

// 播放地址
func play(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	bid := r.FormValue("bid")
	num := r.FormValue("num")
	p, err := strconv.Atoi(num)
	if err != nil {
		requestError(w)
		return
	}
	ab := GetPlaySrc(bid, p)
	resultSucceed(w, CheckURL(ab), 1)
}

func requestFailed(w http.ResponseWriter) {
	var rq RequestResult
	rq.Data = nil
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
	rq.Data = nil
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
	rq.Data = bangumi
	rq.Total = total
	rq.Code = 0
	rq.Message = "succeed"
	jsonStr, err := json.Marshal(rq)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	fmt.Fprintf(w, string(jsonStr))
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

