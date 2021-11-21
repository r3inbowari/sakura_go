package Sakura

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"sync"
)

// RedirectURL 重定向地址变换
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

// GetMD5 md5
func GetMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// SearchDoc 搜索
var SearchDoc sync.Map

func GetSearchResult(url string) *goquery.Document {
	md5Stream := GetMD5(url)
	if v, ok := SearchDoc.Load(md5Stream); ok {
		Log.Info("[Cache] found a snapshot -> " + url)
		return v.(*goquery.Document)
	}
	doc, _ := goquery.NewDocument(url)
	SearchDoc.Store(md5Stream, doc)
	Log.Info("[Cache] snapshot saved -> " + url)
	return doc
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func CreateUUID() string {
	u1 := uuid.NewV4()
	return u1.String()
}

// CreateMD5 md
func CreateMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// JsonBind bind
func JsonBind(ptr interface{}, rq *http.Request) error {
	if rq.Body != nil {
		defer rq.Body.Close()
		err := json.NewDecoder(rq.Body).Decode(ptr)
		if err != nil && err != io.EOF {
			return err
		}
		return nil
	} else {
		return errors.New("empty request body")
	}
}

// CheckURL 重定向检查
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

func ResponseCommon(w http.ResponseWriter, data interface{}, msg string, total int, tag int, code int) {
	var rq RequestResult
	rq.Data = data
	rq.Total = total
	rq.Code = code
	rq.Message = msg
	jsonStr, err := json.Marshal(rq)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err.Error()}).Error("response error")
	}
	w.WriteHeader(tag)
	_, _ = fmt.Fprintf(w, string(jsonStr))
}