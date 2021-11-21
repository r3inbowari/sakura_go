package Sakura

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Player struct {
	ID       string `json:"id"`
	NextLink string `json:"url_next"`
	Link     string `json:"url"`
	Sid      int64  `json:"sid"`
	Nid      int64  `json:"nid"`
	From     string `json:"from"`
}

type Bangumi struct {
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

var DMSakuraHost = "https://www.sakuradm.tv/"

type Detail struct {
	SrcList []int  `json:"src_list"`
	Desc    string `json:"desc"`
}

func DMSakuraDetail(id string) *Detail {
	doc, err := GetSearchResult(fmt.Sprintf("%svoddetail/%s.html", DMSakuraHost, id))
	if err != nil {
		Log.Error("[API] error request")
		return nil
	}
	var ret Detail
	doc.Find(".movurl").Each(func(i int, selection *goquery.Selection) {
		ret.SrcList = append(ret.SrcList, selection.Find("li").Size())
	})
	ret.Desc = doc.Find(".fire.l").Find(".info").Find("p").First().Text()
	return &ret
}

// DMSakuraSearch 关键字搜索0
func DMSakuraSearch(keyword string, page int) ([]Bangumi, int) {
	doc, err := GetSearchResult(fmt.Sprintf("%svodsearch/page/%d/wd/%s.html", DMSakuraHost, page, keyword))
	if err != nil {
		Log.Error("[API] error request")
		return nil, 0
	}
	a := doc.Find(".fire").Find("li")
	var sbs []Bangumi
	a.Each(func(i int, selection *goquery.Selection) {
		var sb Bangumi
		sb.Title = selection.Find("img").AttrOr("alt", "null")
		sb.Cover = selection.Find("img").AttrOr("src", "null")
		sb.Detail = selection.Find("p").Text()
		sb.DetailURL = selection.Find("a").AttrOr("href", "null")
		sb.Update = selection.Find("span").Last().Text()
		arr := strings.Split(sb.Update, "　")
		sb.Update = arr[0]
		if len(arr) > 1 {
			sp := strings.TrimLeft(arr[1], "类型：")
			sp = strings.TrimRight(sp, " ")
			sb.Type = strings.Split(sp, " ")
		}
		sb.ID = strings.TrimLeft(sb.DetailURL, "/voddetail/")
		sb.ID = strings.TrimRight(sb.ID, ".html")
		sbs = append(sbs, sb)
	})
	strCount := doc.Find(".tame em").Text()
	count, err := strconv.Atoi(strCount)
	if err == nil {
		//return sbs, (count + 9) / 10
		return sbs, count
	}
	return sbs, 1
}

// DMSakuraWeeks 周内容
func DMSakuraWeeks(w int) []Bangumi {
	update := LatestHome.Find(".side .bg").Find(".tists").Find("ul")
	var sbs []Bangumi
	update.Each(func(i int, selection *goquery.Selection) {
		if w == i {
			selection.Find("li").Each(func(j int, sel *goquery.Selection) {
				var sb Bangumi
				sb.DetailURL = sel.Find("a").Last().AttrOr("href", "null")
				sb.Title = sel.Find("a").Last().AttrOr("title", "null")
				sb.Update = sel.Find("span").Text()

				sb.ID = strings.TrimLeft(sb.DetailURL, "/voddetail/")
				sb.ID = strings.TrimRight(sb.ID, ".html")

				sbs = append(sbs, sb)
			})
		}
	})
	return sbs
}

// DMSakuraNew 最新内容
func DMSakuraNew() []Bangumi {
	update := LatestHome.Find(".firs.l").Find(".imgs").First().Find("li")
	var sbs []Bangumi
	update.Each(func(i int, selection *goquery.Selection) {
		var sb Bangumi
		sb.Title = selection.Find("img").AttrOr("alt", "null")
		sb.Cover = selection.Find("img").AttrOr("src", "null")
		sb.Update = selection.Find("p").Last().Text()
		sb.DetailURL = selection.Find("a").AttrOr("href", "null")
		sb.ID = strings.TrimLeft(sb.DetailURL, "/voddetail/")
		sb.ID = strings.TrimRight(sb.ID, ".html")
		sbs = append(sbs, sb)
	})
	return sbs
}

// DMSakuraRankUpdate 排行榜
func DMSakuraRankUpdate() []Bangumi {
	update := LatestHome.Find(".side .bg").Last().Find(".pics").Find("li")
	var sbs []Bangumi
	update.Each(func(i int, selection *goquery.Selection) {
		var sb Bangumi
		u := selection.Find("img")
		sb.Title = u.AttrOr("alt", "null")
		sb.Cover = u.AttrOr("src", "null")
		sb.DetailURL = selection.Find("a").AttrOr("href", "null")
		sb.ID = strings.TrimLeft(sb.DetailURL, "/voddetail/")
		sb.ID = strings.TrimRight(sb.ID, ".html")
		sb.Detail = selection.Find("p").Last().Text()
		sbs = append(sbs, sb)
	})
	return sbs
}

func DMSakuraHomepageSnapshotCron() {
	cacheStart = time.Now()
	err := QueryGet(DMSakuraHost, func(doc *goquery.Document) error {
		LatestHome = doc
		return nil
	})
	if err != nil {
		Log.Warn("[Cache] cache pull failed")
	}
	spend := time.Now().UnixNano() - cacheStart.UnixNano()
	Log.Info("[Cache] Homepage snapshot generated | " + strconv.Itoa(int(spend)/1e6) + "ms")
}

func (p *Player) DMSakuraGetPlayer() error {
	get, err := RDB.Get(GetMD5(p.ID + strconv.Itoa(int(p.Sid)) + strconv.Itoa(int(p.Nid))))
	if err == nil {
		err = json.Unmarshal([]byte(get), p)
		return err
	}

	return QueryGet(fmt.Sprintf("%svodplay/%s-%d-%d.html", DMSakuraHost, p.ID, p.Sid, p.Nid), func(doc *goquery.Document) error {
		origin := ""
		doc.Find(".player").Each(func(i int, s *goquery.Selection) {
			origin = s.Text()
		})
		if !strings.Contains(origin, "player_aaaa") {
			return errors.New("502")
		}
		vm := otto.New()
		_, err := vm.Run(origin)
		if err != nil {
			return errors.New("502")
		}
		player, _ := vm.Get("player_aaaa")
		if player.IsDefined() {
			value, _ := player.Object().Get("url_next")
			p.NextLink = value.String()
			value, _ = player.Object().Get("url")
			p.Link = value.String()
			value, _ = player.Object().Get("sid")
			p.Sid, _ = value.ToInteger()
			value, _ = player.Object().Get("nid")
			p.Nid, _ = value.ToInteger()
			value, _ = player.Object().Get("from")
			p.From = value.String()
			value, _ = player.Object().Get("id")
			p.From = value.String()
		}

		marshal, err := json.Marshal(p)
		if err != nil {
			return err
		}
		RDB.SetEx(GetMD5(p.ID+strconv.Itoa(int(p.Sid))+strconv.Itoa(int(p.Nid))), string(marshal), time.Hour)
		return nil
	})
}

func QueryGet(url string, f func(doc *goquery.Document) error) error {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		Log.Error("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		Log.Error("create doc failed")
		return err
	}
	return f(doc)
}
