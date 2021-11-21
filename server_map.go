package Sakura

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func HandleGetPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var p Player
	p.ID = vars["id"]
	p.Sid, _ = strconv.ParseInt(vars["sid"], 10, 64)
	p.Nid, _ = strconv.ParseInt(vars["nid"], 10, 64)
	err := p.DMSakuraGetPlayer()
	if err != nil {
		ResponseCommon(w, "e", "error", 1, http.StatusInternalServerError, -1)
		return
	}
	ResponseCommon(w, p, "ok", 1, http.StatusOK, 0)
}

// HandleRank 排行榜
func HandleRank(w http.ResponseWriter, r *http.Request) {
	ru := DMSakuraRankUpdate()
	ResponseCommon(w, ru, "ok", 1, http.StatusOK, 0)
}

// HandleNew  排行榜
func HandleNew(w http.ResponseWriter, r *http.Request) {
	ru := DMSakuraNew()
	ResponseCommon(w, ru, "ok", 1, http.StatusOK, 0)
}

// HandleWeeks 每周内容
func HandleWeeks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dayIndex := vars["index"]
	c, err := strconv.Atoi(dayIndex)
	if c < 7 && err == nil {
		ru := DMSakuraWeeks(c)
		ResponseCommon(w, ru, "ok", 1, http.StatusOK, 0)
	} else {
		Log.Warn("[API] apply error param")
		ResponseCommon(w, "apply error param", "ok", 1, http.StatusInternalServerError, -1)
	}
}

// HandleSearch 关键字搜索
// /vodsearch/page/1/wd/{keyword}.html
func HandleSearch(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	keyword := r.FormValue("keyword")
	page := r.FormValue("page")
	p, err := strconv.Atoi(page)
	if err == nil {
		ab, ba := DMSakuraSearch(keyword, p)
		ResponseCommon(w, ab, "ok", ba, http.StatusOK, 0)
	} else {
		Log.Warn("[API] apply error param")
		ResponseCommon(w, "apply error param", "ok", 1, http.StatusOK, 0)
	}
}