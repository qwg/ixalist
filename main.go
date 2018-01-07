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
	fmt.Printf("カード番号,名前,コスト,槍\n")
	doc.Find(".parameta_area").Each(func(_ int, s *goquery.Selection) {
		yari := "E"
		if s.Find(".yari.lv_d").Size() == 1 {
			yari = "D"
		} else if s.Find(".yari.lv_c").Size() == 1 {
			yari = "C"
		} else if s.Find(".yari.lv_b").Size() == 1 {
			yari = "B"
		} else if s.Find(".yari.lv_a").Size() == 1 {
			yari = "A"
		} else if s.Find(".yari.lv_s").Size() == 1 {
			yari = "S"
		} else if s.Find(".yari.lv_ss").Size() == 1 {
			yari = "SS"
		} else if s.Find(".yari.lv_sss").Size() == 1 {
			yari = "SSS"
		}
		fmt.Printf("%s,%s,%s,%s\n",
			s.Find(".ig_card_cardno").Text(),
			s.Find(".ig_card_name").Text(),
			s.Find(".ig_card_cost").Text(), yari)
	})
}

func tekisei(s *goquery.Selection, kind string) string {
	level := []string{"sss", "ss", "s", "a", "b", "c", "d", "e"}
	for _, l := range level {
		if s.Find("."+kind+".lv_"+l) == 1 {
			return strings.ToUpper(l)
		}
	}
	return ""
}
