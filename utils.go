package Sakura

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"sync"
)

/*
 * 重定向地址变换
 * @param baseURL origin visit url
 */
func RedirectURL(baseURL string) string {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for {
		res, err := client.Get(baseURL)
		if err != nil {
			return baseURL
		}
		if res.StatusCode != 302 {
			return baseURL
		}
		baseURL = res.Header.Get("Location")
	}
}

/**
 * MD5生成
 */
func GetMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

/**
 * 搜索结果带缓存
 */
var SearchDoc sync.Map

func GetSearchResult(url string) *goquery.Document {
	md5Stream := GetMD5(url)
	if v, ok := SearchDoc.Load(md5Stream); ok {
		INFO("提取地址: " + url)
		return v.(*goquery.Document)
	}
	doc, _ := goquery.NewDocument(url)
	SearchDoc.Store(md5Stream, doc)
	INFO("缓存地址: " + url)
	return doc
}

