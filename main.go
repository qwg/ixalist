package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	var filename string
	flag.StringVar(&filename, "string", "C:\\data\\a.html", "ixa listfile")
	fileInfos, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	stringReader := strings.NewReader(string(fileInfos))
	doc, err := goquery.NewDocumentFromReader(stringReader)
	if err != nil {
		fmt.Print("url scarapping failed")
	}
	doc.Find(".line-number").Each(func(_ int, s *goquery.Selection) {
		fmt.Println(s)
	})
}
