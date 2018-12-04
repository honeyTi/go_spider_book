package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

func HttpGet(url string) (result string, err error)  {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err =err1
		return
	}
	defer resp.Body.Close()
	// 读取网页内容
	buf := make([]byte,1024*4)
	for {
		n,_ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		result += string(buf[:n]) // 累加读取的内容
	}
	return
}

// 爬取章节内容
func GetBookContent(url string) (bookIndex, bookContent string, err error)  {
	result, err := HttpGet(url)
	if err != nil {
		fmt.Println("get book content HttpGet() method error : ", err)
		return
	}
	// 取书名
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

func httpToDatabase(bookName, bookIndex, bookContent string)  {
	data := make(url.Values)
	data["book_name"] = []string{bookName}
	data["book_index"] = []string{bookIndex}
	data["book_content"] = []string{bookContent}
	resp, err :=  http.PostForm("http://127.0.0.1:8090/addBook25", data)
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

// 获取爬取页面内容
func DoWork(url string)  {
	//开始获取书籍简绍主页内容
	result, err := HttpGet(url)
	if err != nil  {
		fmt.Println("httpGet() method error : ", err)
		return
	}

	regx := regexp.MustCompile(`<td><a href="(?s:(.*?))"`)

	if regx == nil {
		fmt.Println("book-目录， 未匹配到相关内容")
		return
	}

	indexUrl := regx.FindAllStringSubmatch(result, -1)

	//名称
	book_index := make([]string, 0)
	book_content := make([]string, 0)
	// 取出目录的网址拼接在一起
	for _, data := range indexUrl {
		indexList := "https://www.dushu.com" + data[1]
		index, content, err := GetBookContent (indexList)
		if err != nil {
			fmt.Println("完成")
			break
		}
		book_index = append(book_index, index)
		book_content = append(book_content, content)
	}

	for i:= 0;i< len(book_index) ;i++  {
		 httpToDatabase("史记", book_index[i], book_content[i])
	}
}

func main() {
	var httpString string
	fmt.Println("请输入爬取书籍的页面网址：")
	fmt.Scan(&httpString)
	fmt.Println(httpString)

	DoWork(httpString)
}
