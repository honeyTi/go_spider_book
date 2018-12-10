package content_map

import (
	"book_spider_go/burnbook.com/book_urls"
	"book_spider_go/burnbook.com/html_get"
	"fmt"
	"regexp"
)

func GetBookContent(book_url, key string) (bookName, bookAuthor, bookDesc string, err error) {
	result, err1 := html_get.HttpGet(book_url)
	if err1 != nil {
		fmt.Println("GetBookContent method err :" , err1)
		err = err1
		return
	}
	// 取书名
	regBookName := regexp.MustCompile(`<h4>(?s:(.*?))</h4>`+`
                    <hr>`)
	if regBookName == nil {
		fmt.Println("book index regx err")
		return
	}
	bookNameAll := regBookName.FindAllStringSubmatch(result, -1)
	for _, data := range bookNameAll {
		bookName = data[1]
		break
	}
	// 取作者名
	regBookAuthor := regexp.MustCompile(`class="text-muted">(?s:(.*?))</a>`)
	if regBookAuthor == nil {
		fmt.Println("无法获取作者名字")
		return
	}
	bookAuthorAll := regBookAuthor.FindAllStringSubmatch(result, -1)
	for _, data := range bookAuthorAll {
		bookAuthor = data[1]
		break
	}
	// 取作品描述
	regBookDesc := regexp.MustCompile(`<div class="detail">(?s:(.*?))</div>`)
	if regBookDesc == nil {
		fmt.Println("无法获取作品描述")
		return
	}
	bookDescAll := regBookDesc.FindAllStringSubmatch(result, -1)
	for _, data := range bookDescAll {
		bookDesc = data[1]
		break
	}
	// 获取书籍的章节
	regBookUrl := regexp.MustCompile(`<a class="card-link" href="(?s:(.*?))">`)
	if regBookUrl == nil {
		fmt.Println("获取书籍章节内容失败")
		return
	}
	regBookUrls := regBookUrl.FindAllStringSubmatch(result,-1)
	for _, data := range regBookUrls {
		urls := "https://burnook.com" + data[1]
		book_urls.IndexGet(urls, key)
	}
	return
}