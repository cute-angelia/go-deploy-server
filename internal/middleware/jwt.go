package middleware

import (
	"fmt"
	"github.com/cute-angelia/go-utils/cache/bunt"
	"github.com/cute-angelia/go-utils/http/api"
	"github.com/cute-angelia/go-utils/http/jwt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// ===== 请勿删除， 可支持传入配置 =====
var allowLoginPaths = []string{
	"/",
	"/favicon.ico",
}

/**
设置头部信息， 包括未登陆的
*/
func setHeaderInfo(r *http.Request, secret string, jwtToken string) {
	ijwt := jwt.NewJwt(secret)
	if jwtDecode, err := ijwt.Decode(jwtToken); err == nil {
		uid, _ := jwtDecode.Get("uid")
		openid, _ := jwtDecode.Get("openid")
		cid, _ := jwtDecode.Get("cid")
		nickname, _ := jwtDecode.Get("nickname")

		head, _ := jwtDecode.Get("head")
		appid, _ := jwtDecode.Get("appid")

		if uid != nil {
			r.Header.Set("jwt_uid", fmt.Sprintf("%v", uid))
		}
		if openid != nil {
			r.Header.Set("jwt_openid", fmt.Sprintf("%v", openid))
		}
		if head != nil {
			r.Header.Set("jwt_head", fmt.Sprintf("%v", head))
		}
		if appid != nil {
			r.Header.Set("jwt_appid", fmt.Sprintf("%v", appid))
		}

		icid := fmt.Sprintf("%v", cid)
		if cid != nil && icid != "0" {
			r.Header.Set("jwt_cid", fmt.Sprintf("%v", cid))
		} else {
			cid := r.URL.Query().Get("cid")
			r.Header.Set("jwt_cid", fmt.Sprintf("%v", cid))
		}

		if nickname != nil {
			r.Header.Set("jwt_nickname", fmt.Sprintf("%v", nickname))
		}
	} else {
		log.Println("setHeaderInfo", err.Error())
	}
}

func Jwt(allowList []string, secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			// 提取 Token
			authToken := r.Header.Get("Authorization")
			jwtToken := ""

			if len(authToken) > 0 && strings.Contains(authToken, "Bearer") {
				allToken := strings.Split(authToken, "Bearer ")
				if len(allToken) == 2 {
					jwtToken = allToken[1]
				}
			}
			r.Header.Set("jwt_token", jwtToken)
			if len(jwtToken) > 0 {
				setHeaderInfo(r, secret, jwtToken)
			}
			// 提取 Token end

			// check logout
			isLogOut := bunt.Get("cache", jwtToken)
			if isLogOut == "true" {
				log.Println("例外匹配过滤JWT ×：", -999, r.URL.Path, jwtToken)
				api.Error(w, r, nil, "登录已过期, 请重新登录", -999)
				return
			}

			// 白名单 例外直接过
			if len(allowList) > 0 {
				for _, vallow := range allowList {
					allowLoginPaths = append(allowLoginPaths, vallow)
				}
			}
			for _, v := range allowLoginPaths {
				// 正则匹配
				if strings.Contains(v, "*") {
					r1 := regexp.MustCompile(v)
					if r1.MatchString(r.URL.Path) {
						next.ServeHTTP(w, r)
						return
					}
				}
				// 路径匹配
				if r.URL.Path == v {
					next.ServeHTTP(w, r)
					return
				}
			}

			// 初步校验
			if len(jwtToken) < 30 {
				log.Println("例外匹配过滤JWT ×：", -999, r.URL.Path, jwtToken)
				api.Error(w, r, nil, "登录已过期, 请重新登录", -999)
				return
			}

			ijwt := jwt.NewJwt(secret)
			if _, err := ijwt.Decode(jwtToken); err != nil {
				w.WriteHeader(203)
				log.Println("例外匹配过滤JWT ×：", -999, r.URL.Path, jwtToken)
				api.Error(w, r, nil, "登录已过期, 请重新登录", -999)
				return
			} else {
				//setHeaderInfo(r, jwtDecode)
			}
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
