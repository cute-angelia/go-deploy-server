package middleware

import (
	"fmt"
	"github.com/cute-angelia/go-utils/cache/bunt"
	"github.com/cute-angelia/go-utils/http/api"
	"github.com/cute-angelia/go-utils/utils/encrypt/wechat"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var allowPathsSign = []string{
	"/sso/auth/wxLogin",
	"/sso/auth/wxLoginDo",
	"/sso/auth/wxLoginGz",
	"/open/record/receive",
	"/sso/auth/getMiniToken",
	"/api-game-reward/customer/verification",
	"/api-game-reward/user/complete",
}

func SignPass(allowList []string, apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			signIn := r.URL.Query().Get("sign")

			// 使用传入配置
			if len(allowList) > 0 {
				allowPathsSign = allowList
			}

			// 例外直接过
			for _, v := range allowPathsSign {
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

			// debug
			if r.URL.Query().Get("debug") == "ins" {
				next.ServeHTTP(w, r)
				return
			}

			sign := wechat.NewSign(apiKey)
			signOut := sign.Signature(r.URL.Query())
			if strings.ToLower(signIn) != strings.ToLower(signOut) {
				log.Printf("signkey:%s, query:%s", sign.GetApiKey(), r.URL.Query())

				w.WriteHeader(203)
				api.Error(w, r, nil, "签名失败", -1)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func SignPre(allowList []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			// 使用传入配置
			if len(allowList) > 0 {
				allowPathsSign = allowList
			}

			// 例外直接过
			for _, v := range allowPathsSign {
				// 正则匹配
				if strings.Contains(v, "*") {
					r1 := regexp.MustCompile(v)

					log.Println("例外匹配过滤SIGN:", r1.MatchString(r.URL.Path), v, r.URL.Path)
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

			// 127.0.0.1
			//if strings.Contains(r.RemoteAddr, "127.0.0.1") {
			//	next.ServeHTTP(w, r)
			//	return
			//}

			// debug
			if r.URL.Query().Get("debug") == "ins" {
				next.ServeHTTP(w, r)
				return
			}

			// 过期判断
			timestamp := r.URL.Query().Get("nonce_time")
			t := int(time.Now().Unix())
			timestampInt, _ := strconv.Atoi(timestamp)
			if t-timestampInt > 500 || timestampInt-t > 360 {
				w.WriteHeader(203)
				log.Println("now -> timestampInt", t, timestampInt)
				msg := fmt.Sprintf("now %d -> timestampInt %d", t, timestampInt)
				api.Error(w, r, nil, "请求已过期"+msg, -1)
				return
			}
			// 过期判断 end

			// nonce 重复请求判断
			nonce := r.URL.Query().Get("nonce_str")

			opt := bunt.NewLockerOpt(bunt.WithToday(false))

			log.Println("SIGNPRE_REPEAT_" + nonce + timestamp)

			if !bunt.IsNotLockedInLimit("middleware", "SIGNPRE_REPEAT_"+nonce+timestamp, time.Minute*6, opt) {
				w.WriteHeader(203)
				api.Error(w, r, nil, "请勿重复请求，当前请求已过期", -1)
				return
			}

			// nonce 重复请求判断 end

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
