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
	find := doc.Find("#text")
	html, _ := find.Html()
	html = strings.Replace(html, "</p>", "\n", -1)
	html = strings.Replace(html, "<p>", "\t", -1)
	title := "\t" + doc.Find(".title").Text() + "\n"
	nextURL, is := doc.Find("#next2").Attr("href")

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
	file.WriteAt([]byte("\t"+title), seek)
	seek, err = file.Seek(0, 2)
	if err != nil {
		fmt.Println(err.Error())
	}
	file.WriteAt([]byte("\t"+html), seek)
	nextURL = nextURL
	if len(nextURL) != 0 {
		GetRes(nextURL, bookName)
	}
}
