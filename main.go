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
		b := NewBusho(s)
		fmt.Println(b)
		/*

			fmt.Printf("%s,%s,%s,%.1f,%s,%s,%s,%s,%s,%s,%s,%s\n",
				s.Find(".ig_card_status_att").Text(),
				s.Find(".ig_card_status_def").Text(),
				s.Find(".ig_card_status_int").Text(),
				strings.Split(s.Find(".commandsol_no").Text(), "/")[1])
		*/
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

type attackKind struct {
	tekisei string
	hosei   int
}

type hei struct {
	name string
	kind string
}

/*
const heilist []hei = {
	{"長槍", "yari"}
}
*/

//Busho ...
type Busho struct {
	no       string
	name     string
	cost     float64
	shubetsu string
	att      int
	def      int
	skill    float64
	comno    int
	yari     attackKind
	kiba     attackKind
	yumi     attackKind
	heiki    attackKind
}

//NewBusho ...
func NewBusho(s *goquery.Selection) *Busho {
	var Busho Busho
	Busho.init(s)
	return &Busho
}

func (b *Busho) init(s *goquery.Selection) {
	if len(s.Find(".ig_card_cost").Text()) == 0 {
		b.cost = -1
	} else {
		b.cost, _ = strconv.ParseFloat(s.Find(".ig_card_cost").Text(), 64)
	}
	b.no = s.Find(".ig_card_cardno").Text()
	b.name = s.Find(".ig_card_name").Text()
	b.shubetsu = shubetsu(s)
	b.att, _ = strconv.Atoi(s.Find(".ig_card_status_att").Text())
	b.def, _ = strconv.Atoi(s.Find(".ig_card_status_def").Text())
	b.skill, _ = strconv.ParseFloat(s.Find(".ig_card_status_int").Text(), 64)
	b.comno, _ = strconv.Atoi(strings.Split(s.Find(".commandsol_no").Text(), "/")[1])
	tekisei(s, "yari", &b.yari)
	tekisei(s, "kiba", &b.kiba)
	tekisei(s, "yumi", &b.yumi)
	tekisei(s, "heiki", &b.heiki)
}

func tekisei(s *goquery.Selection, kind string, k *attackKind) {
	level := []string{"sss", "ss", "s", "a", "b", "c", "d", "e", "f"}
	for i, l := range level {
		if s.Find("."+kind+".lv_"+l).Size() == 1 {
			k.tekisei = strings.ToUpper(l)
			k.hosei = 120 - i*5
			return
		}
	}
}
