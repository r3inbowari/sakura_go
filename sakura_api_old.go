package Sakura

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
	"time"
)

// 排行榜
func RankUpdate() []Bangumi {
	doc, _ := goquery.NewDocument("http://www.yhdm.tv")
	update := doc.Find(".side").Find(".pics").Find("li")
	var sbs []Bangumi
	update.Each(func(i int, selection *goquery.Selection) {
		var sb Bangumi
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

// 最新内容
func LastUpdate() []Bangumi {
	// doc, _ = goquery.NewDocument("http://www.yhdm.tv")
	start := time.Now().UnixNano()
	update := LatestHome.Find(".firs").Find(".img").First().Find("li")
	var sbs []Bangumi
	update.Each(func(i int, selection *goquery.Selection) {
		var sb Bangumi
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

// 周内容
func WeekUpdate(w int) []Bangumi {
	doc, _ := goquery.NewDocument("http://www.yhdm.tv")
	update := doc.Find(".side").Find(".tlist").Find("ul")
	var sbs []Bangumi
	update.Each(func(i int, selection *goquery.Selection) {
		if w == i {
			selection.Find("li").Each(func(j int, sel *goquery.Selection) {
				var sb Bangumi
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

// 关键字搜索0
func SearchBangumi(keyword string, page int) ([]Bangumi, int) {
	doc, _ := GetSearchResult("http://www.yhdm.tv/search/" + keyword + "/?page=" + strconv.Itoa(page))
	a := doc.Find(".lpic").Find("li")
	var sbs []Bangumi
	a.Each(func(i int, selection *goquery.Selection) {
		var sb Bangumi
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

// 关键字搜索1
func SearchBangumi1(keyword string, page int) ([]Bangumi, int) {
	doc, _ := GetSearchResult("http://www.yhdm.tv/search/" + keyword + "/?page=" + strconv.Itoa(page))
	a := doc.Find(".lpic").Find("li")
	var sbs []Bangumi
	a.Each(func(i int, selection *goquery.Selection) {
		var sb Bangumi
		sb.Title = selection.Find("img").AttrOr("alt", "null")
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

// 详细页
func DetailBangumi(id string) Bangumi {
	doc, _ := GetSearchResult("http://www.yhdm.tv/show/" + id + ".html")
	// doc, _ := goquery.NewDocument("http://www.yhdm.tv/show/" + id + ".html")
	pInfo := doc.Find(".sinfo").Find("p")
	var sb Bangumi
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

// 播放地址
func GetPlaySrc(id string, num int) string {
	url := "http://www.yhdm.tv/v/" + id + "-" + strconv.Itoa(num) + ".html"
	doc, _ := GetSearchResult(url)
	if doc == nil {
		SearchDoc.Delete(url)
		return ""
	}
	playSrc := doc.Find(".bofang").Find("div").AttrOr("data-vid", "null")

	defer func() {
		if e := recover(); e != nil {
			Log.Warn("[INFO] 地址解析失败: url")
		}
	}()

	inx := strings.LastIndex(playSrc, "$")
	return playSrc[:inx]
}

func (sb *Bangumi) BangumiDetail() error {
	doc, _ := GetSearchResult("http://www.yhdm.tv/show/" + sb.ID + ".html")
	pInfo := doc.Find(".sinfo").Find("p")
	if pInfo.Size() > 1 {
		sb.Alias = strings.TrimLeft(pInfo.First().Text(), "别名:")
	}
	sb.Detail = doc.Find(".info").Text()
	sb.Count = doc.Find(".movurl").Find("li").Size()

	sb.Title = doc.Find(".thumb").Find("a").AttrOr("title", "null")
	sb.Cover = doc.Find("img").AttrOr("src", "null")

	sb.DetailURL = "/show/" + sb.ID + ".html"
	sb.Update = pInfo.Last().Text()

	doc.Find(".sinfo").Find("span").First().Next().Next().Find("a").Each(func(i int, selection *goquery.Selection) {
		sb.Type = append(sb.Type, selection.AttrOr("href", "null"))
	})
	return nil
}
