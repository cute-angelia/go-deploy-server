package controller

import (
	"encoding/json"
	"github.com/cute-angelia/go-utils/http/api"
	"github.com/go-chi/chi/v5"
	"go-deploy/config"
	"net/http"
)

type List struct {
}

func (rs List) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", rs.list)
	return r
}

// 返回结构，过滤敏感信息
type ListResp struct {
	GroupId string `json:"groupid"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Node    []struct {
		Alias  string `json:"alias"`
		Online bool   `json:"online"`
	} `json:"node"`
}

func (rs List) list(w http.ResponseWriter, r *http.Request) {
	resps := []ListResp{}
	temp, _ := json.Marshal(config.C.Apps)
	json.Unmarshal(temp, &resps)

	api.Success(w, r, resps, "success")
	return
}
