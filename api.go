package Sakura

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
)

/**
 * 搜索
 */
func SearchBangumi(keyword string, page int) ([]SakuraBangumi, int) {
	doc := GetSearchResult("http://www.yhdm.tv/search/" + keyword + "/?page=" + strconv.Itoa(page))
	a := doc.Find(".lpic").Find("li")
	var sbs []SakuraBangumi
	a.Each(func(i int, selection *goquery.Selection) {
		var sb SakuraBangumi
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

/**
 * 最新更新
 * 五分钟内不更新
 */
func LatestUpdate() []SakuraBangumi {
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
	return sbs
}

/**
 * 排行榜
 */
func RankUpdate() []SakuraBangumi {
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

/**
 * 每周更新
 */
func WeekUpdate(w int) []SakuraBangumi {
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

func (sb *SakuraBangumi) BangumiDetail() error {
	doc := GetSearchResult("http://www.yhdm.tv/show/" + sb.ID + ".html")
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

func (sb *SakuraBangumi) getPlaySrc(num int) string {
	doc := GetSearchResult("http://www.yhdm.tv/v/" + sb.ID + "-" + strconv.Itoa(num) + ".html")
	playSrc := doc.Find(".bofang").Find("div").AttrOr("data-vid", "null")

	defer func() {
		if e := recover(); e != nil {
			log.Println("[INFO] 地址解析失败: url")
		}
	}()

	inx := strings.LastIndex(playSrc, "$")
	return playSrc[:inx]
}

