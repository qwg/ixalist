package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
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
	fmt.Printf("カード番号,名前,種別,コスト,槍,馬,弓,器,攻,防,兵,指揮力,攻(槍)\n")
	doc.Find(".card_detail_area").Each(func(_ int, s *goquery.Selection) {
		if len(s.Find(".ig_card_cost").Text()) == 0 {
			return
		}
		cost, _ := strconv.ParseFloat(s.Find(".ig_card_cost").Text(), 64)
		yari := tekisei(s, "yari")
		kiba := tekisei(s, "kiba")
		yumi := tekisei(s, "yumi")
		heiki := tekisei(s, "heiki")
		shubetsu := shubetsu(s)
		fmt.Printf("%s,%s,%s,%.1f,%s,%s,%s,%s,%s,%s,%s,%s\n",
			s.Find(".ig_card_cardno").Text(),
			s.Find(".ig_card_name").Text(),
			shubetsu,
			cost, yari, kiba, yumi, heiki,
			s.Find(".ig_card_status_att").Text(),
			s.Find(".ig_card_status_def").Text(),
			s.Find(".ig_card_status_int").Text(),
			strings.Split(s.Find(".commandsol_no").Text(), "/")[1])
	})
}

func shubetsu(s *goquery.Selection) string {
	kind := []string{"0", "将", "剣", "忍", "文", "姫", "6", "7"}
	for i, t := range kind {
		if s.Find(".jobtype_"+strconv.Itoa(i)).Size() == 1 {
			return t
		}
	}
	return "?"
}

func tekisei(s *goquery.Selection, kind string) string {
	level := []string{"sss", "ss", "s", "a", "b", "c", "d", "e", "f"}
	for _, l := range level {
		if s.Find("."+kind+".lv_"+l).Size() == 1 {
			return strings.ToUpper(l)
		}
	}
	return ""
}
