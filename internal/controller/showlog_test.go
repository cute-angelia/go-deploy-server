package controller

import (
	"log"
	"regexp"
	"testing"
)

func TestLog(t *testing.T) {
	str := `------------------------------------------------------------------------
r756 | wangjianzhong | 2020-12-17 09:58:52 +0800 (▒▒▒▒, 17 12▒▒ 2020) | 1 line


------------------------------------------------------------------------
r755 | wangjianzhong | 2020-10-12 13:45:06 +0800 (▒▒һ, 12 10▒▒ 2020) | 1 line


------------------------------------------------------------------------
r754 | wangjianzhong | 2020-09-02 17:13:54 +0800 (▒▒▒▒, 02 9▒▒ 2020) | 1 line


------------------------------------------------------------------------
r753 | wangjianzhong | 2020-09-02 17:12:40 +0800 (▒▒▒▒, 02 9▒▒ 2020) | 1 line
`
	svnlogRegex := regexp.MustCompile(`r(\d+) \| (\w+) \| (.*) \+0800(?:.*)\n\n(.*)\n--`)
	match := svnlogRegex.FindAllStringSubmatch(str, -1)
	for _, item := range match {
		for _, s := range item {
			log.Println(s)
		}
	}
}
