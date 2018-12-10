package book_urls

import (
	"book_spider_go/burnbook.com/html_get"
	"fmt"
	"regexp"
)

func IndexGet(url, key string) {
	result, err := html_get.HttpGet(url)
	if err != nil {
		fmt.Println("get index content error :", err)
		return
	}
	// 获取书籍名称
	regex := regexp.MustCompile(`<a href="`+key+`">(?s:(.*?))</a>`)

	if regex == nil {
		fmt.Println("完成章节内容获取")
	}
	regexName := regex.FindAllStringSubmatch(result, -1)

	for _, data := range regexName {
		regex_name := data[1]
	}
	// 获取书籍章节名称
	regexIndex := regexp.MustCompile(`<a href="`+key+`">(?s:(.*?))</a>`)

	if regexIndex == nil {
		fmt.Println("完成章节内容获取")
	}
	regexName := regex.FindAllStringSubmatch(result, -1)

	for _, data := range regexName {
		regex_name := data[1]
	}
}