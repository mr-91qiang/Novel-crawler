package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	var url, bookName string
	flag.StringVar(&url, "u", "", "路径")
	flag.StringVar(&bookName, "n", "", "书名")
	flag.Parse()
	fmt.Println(url, bookName)

	//GetRes("https://www.changyeyuhuo.com/book/21955/3895357.html","长夜余火.text")
	GetRes(url, bookName)
}

func GetRes(url string, bookName string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	find := doc.Find("#content")
	html := find.Text()
	html = strings.Replace(html, "。", "。\n", -1)
	title := doc.Find(".title").Text()
	nextURL, is := doc.Find("#next_url").Attr("href")

	if !is {
		fmt.Println("没有下一页了")
		return
	}
	file, err := os.OpenFile(bookName, os.O_CREATE|os.O_RDWR, 0666)
	defer file.Close()
	seek, err := file.Seek(0, 2)
	if err != nil {
		fmt.Println(err.Error())
	}
	file.WriteAt([]byte(title), seek)
	seek, err = file.Seek(0, 2)
	if err != nil {
		fmt.Println(err.Error())
	}
	file.WriteAt([]byte(html), seek)
	nextURL = "https://www.changyeyuhuo.com" + nextURL
	GetRes(nextURL, bookName)
}
