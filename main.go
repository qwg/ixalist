package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 兵士編成を表示
// ソースを保存
// chromeで表示
// 表示したものを保存
// 保存したファイルを引数で指定
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
	fmt.Printf("カード番号,名前,コスト\n")
	doc.Find(".parameta_area").Each(func(_ int, s *goquery.Selection) {
		fmt.Printf("%s,%s,%s\n",
			s.Find(".ig_card_cardno").Text(),
			s.Find(".ig_card_name").Text(),
			s.Find(".ig_card_cost").Text())
	})
}
