package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()
	// 读取网页内容
	buf := make([]byte, 1024*4)
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		result += string(buf[:n]) // 累加读取的内容
	}
	return
}

// 爬取章节内容
func GetBookContent(url, book string) (booName, bookIndex, bookContent string, err error) {
	result, err := HttpGet(url)
	if err != nil {
		fmt.Println("get book content HttpGet() method error : ", err)
		return
	}
	// 取书名
	regBookName := regexp.MustCompile(`<p class="text-center padding-top"><a href="` + book + `">(?s:(.*?))</a>`)

	if regBookName == nil {
		fmt.Println("book index regx err")
		return
	}
	bookNameAll := regBookName.FindAllStringSubmatch(result, -1)
	for _, data := range bookNameAll {
		booName = data[1]
		break
	}
	// 取目录
	regBookIndex := regexp.MustCompile(`<p class="text-center text-large padding-top">(?s:(.*?))</p>`)

	if regBookIndex == nil {
		fmt.Println("book index regx err")
		return
	}
	bookIndexAll := regBookIndex.FindAllStringSubmatch(result, -1)
	for _, data := range bookIndexAll {
		bookIndex = data[1]
		break
	}
	// 取内容
	regBookContent := regexp.MustCompile(`<div class="content_txt">(?s:(.*?))</div>`)

	if regBookContent == nil {
		fmt.Println("book index regx err")
		return
	}
	bookContentAll := regBookContent.FindAllStringSubmatch(result, -1)
	for _, data := range bookContentAll {
		bookContent = data[1]
		break
	}
	return
}

func httpToDatabase(bookName, bookIndex, bookContent, bookKey string) {
	data := make(url.Values)
	data["book_name"] = []string{bookName}
	data["book_index"] = []string{bookIndex}
	data["book_content"] = []string{bookContent}
	data["book_key"] = []string{bookKey}
	resp, err := http.PostForm("http://127.0.0.1:8090/addBook25", data)
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

// 获取书籍内容
func GetBookDesc(result, url string) {
	var bookName, bookAuthor, bookDesc string

	// 获取书籍名字
	reg_book_name := regexp.MustCompile(`<h1>(?s:(.*?))</h1>`)
	if reg_book_name == nil {
		return
	}
	book_name := reg_book_name.FindAllStringSubmatch(result, -1)
	for _, data := range book_name {
		bookName = data[1]
		break
	}
	// 获取书籍作者
	reg_book_author := regexp.MustCompile(`title="更多同作者相关图书" target="_blank">(?s:(.*?))</a>`)
	if reg_book_author == nil {
		return
	}
	book_author := reg_book_author.FindAllStringSubmatch(result, -1)
	for _, data := range book_author {
		bookAuthor = data[1]
		break
	}
	// 获取书籍描述
	reg_book_desc := regexp.MustCompile(`<div class="text txtsummary">(?s:(.*?))</div>`)
	if reg_book_desc == nil {
		return
	}
	book_desc := reg_book_desc.FindAllStringSubmatch(result, -1)
	for _, data := range book_desc {
		bookDesc = data[1]
		break
	}

	BookDetailToDatabase(bookName, bookAuthor, bookDesc , url)
}

// 获取爬取页面内容
func DoWork(url, book string) {
	//开始获取书籍简绍主页内容
	result, err := HttpGet(url)
	if err != nil {
		fmt.Println("httpGet() method error : ", err)
		return
	}

	GetBookDesc(result, book)

	regx := regexp.MustCompile(`<td><a href="(?s:(.*?))"`)

	if regx == nil {
		fmt.Println("book-目录， 未匹配到相关内容")
		return
	}

	indexUrl := regx.FindAllStringSubmatch(result, -1)

	//名称
	book_name := make([]string, 0)
	book_index := make([]string, 0)
	book_content := make([]string, 0)
	book_key := make([]string, 0)
	// 取出目录的网址拼接在一起
	for _, data := range indexUrl {
		indexList := "https://www.dushu.com" + data[1]
		name, index, content, err := GetBookContent(indexList, book)
		if err != nil {
			fmt.Println("完成")
			break
		}
		book_name = append(book_name, name)
		book_index = append(book_index, index)
		book_content = append(book_content, content)
		book_key = append(book_key, book)
	}

	for i:= 0;i< len(book_index) ;i++  {
		 httpToDatabase(book_name[i], book_index[i], book_content[i], book)
	}
}


// 爬取书籍列表页面
func BookDetail(url string) {
	result, err := HttpGet(url)
	if err != nil {
		fmt.Println("HttpGet method error in BookDetail method :", err)
		return
	}
	//获取book书名---作者----简介----详情页url

	// 书籍url
	book_url_rgx := regexp.MustCompile(`<h3><a href="(?s:(.*?))"`)
	if book_url_rgx == nil {
		return
	}
	urls := book_url_rgx.FindAllStringSubmatch(result, -1)
	for _, data := range urls {
		indexList := "https://www.dushu.com" + data[1]
		// 获取书籍内容
		DoWork(indexList, data[1])
	}
}

// 爬取书籍内容，及链接网址
func BookList(url, endPage string) {
	urls := make([]string, 0)
	page, err := strconv.Atoi(endPage)
	if err != nil {
		fmt.Println("endpage 页码错误")
		return
	}
	if page == 1 {
		urls = append(urls, url)
	} else {
		for i := 1; i <= page; i++ {
			m := strconv.Itoa(i)
			str := "_" + m + ".html"
			urls = append(urls, strings.Replace(url, ".html", str, -1))
		}
	}
	for _, data := range urls {
		 BookDetail(data)
	}
}

func main() {
	var httpString, endPage string
	fmt.Println("请输入爬取书籍的页面首页地址：")
	fmt.Scan(&httpString)
	fmt.Println("爬取页面结束页码（请输入数字, 如果只有一页，请输入1）：")
	fmt.Scan(&endPage)

	BookList(httpString, endPage)
}
