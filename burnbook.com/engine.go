package main

import (
	"book_spider_go/burnbook.com/content-map"
	"book_spider_go/burnbook.com/html_get"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)
// 写入数据库
func BookDetailToDatabase(bookName, bookAuthor, bookDesc, bookUrl string) {
	data := make(url.Values)
	data["book_name"] = []string{bookName}
	data["book_author"] = []string{bookAuthor}
	data["book_desc"] = []string{bookDesc}
	data["book_key"] = []string{bookUrl}
	resp, err := http.PostForm("http://127.0.0.1:8090/addbookdetail", data)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}
// 启动go程序
func DoWork(url string)  {
	result, err := html_get.HttpGet(url)
	if err != nil {
		fmt.Println("DoWork method httpGet error: ", err)
		return
	}
	regx := regexp.MustCompile(`<h5><a href="(?s:(.*?))"`)

	if regx == nil {
		fmt.Println("获取网址出错")
	}

	indexUrl := regx.FindAllStringSubmatch(result,-1)
	// 取出字页面地址
	for _, data := range indexUrl {
		book_url := "https://burnook.com" + data[1]
		bookName, bookAnthor, bookDesc, err := content_map.GetBookContent(book_url, data[1])
		if err != nil {
			fmt.Println("完成")
			break
		}
		fmt.Println(bookName)
		fmt.Println(bookAnthor)
		fmt.Println(bookDesc)
		fmt.Println(data[1])
		//BookDetailToDatabase(bookName, bookAnthor, bookDesc, data[1])
	}
}
// 获取爬取页面信息
func UrlGet(startPage, endPage string) {
	start, _ := strconv.Atoi(startPage)
	end, _ := strconv.Atoi(endPage)
	for i := start; i <= end; i++ {
		page := strconv.Itoa(i)
		urlMap := "https://burnook.com/documents/list?page=" + page + `&cate=15`
		DoWork(urlMap)
	}
}

func main() {
	var startPage, endPage string
	fmt.Println("请输入起始页码：")
	fmt.Scan(&startPage)
	fmt.Println("请输入结束页码：")
	fmt.Scan(&endPage)

	UrlGet(startPage, endPage)
}
