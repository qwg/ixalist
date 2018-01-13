package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// 兵士編成を表示
// ソースを保存
// chromeで表示
// 表示したものを保存
// 保存したファイルを引数で指定
func main() {
	var filename string
	flag.StringVar(&filename, "f", "C:\\data\\b.html", "ixa listfile")
	flag.Parse()
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

	heiList := []Hei{
		{"長槍", "槍", "", 18, 18},
		{"長弓", "弓", "", 17, 19},
		{"精鋭騎馬", "馬", "", 19, 16},
		{"赤備え", "馬", "槍", 23, 20},
		{"武士", "槍", "弓", 22, 22},
		{"弓騎馬", "弓", "馬", 21, 23},
	}

	fmt.Printf(utf82sjis("カード番号,名前,種別,コスト,槍,馬,弓,器,攻,防,スキル,指揮力,攻撃種,攻撃力,防御種,防御力,コスト比\n"))
	doc.Find(".card_detail_area").Each(func(_ int, s *goquery.Selection) {
		b := NewBusho(s, heiList)
		maxAttackName, maxAttackScore := maxAttack(b)
		maxDefName, maxDef, maxDefCost := maxDef(b)
		str := fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
			b.no, b.name, b.shubetsu, b.cost, b.yari.tekisei, b.kiba.tekisei, b.yumi.tekisei, b.heiki.tekisei,
			b.att, b.def, b.skill, b.comno, maxAttackName, maxAttackScore, maxDefName, maxDef, maxDefCost)
		fmt.Print(utf82sjis(str))
		//fmt.Print(str)
		/*
			max_shubetsu, max_cost := max(b)
			str := fmt.Sprintf("%s,%s,%s,%.1f,%s,%s,%s,%s,%d,%d,%.0f,%d,%d,%d,%d,%d,%d,%s,%d\n",
				b.no, b.name, b.shubetsu, b.cost, b.yari.tekisei, b.kiba.tekisei, b.yumi.tekisei, b.heiki.tekisei,
				b.att, b.def, b.skill, b.nagayari.def, b.nagayari.defCost,
				b.nagayumi.def, b.nagayumi.defCost, b.seieikiba.def, b.seieikiba.defCost, max_shubetsu, max_cost)
			fmt.Print(utf82sjis(str))
			//fmt.Print(str)
		*/
	})
}

func maxAttack(b *Busho) (name string, att int) {
	for _, tekisei := range b.tekisei {
		if tekisei.att > att {
			att = tekisei.att
			name = tekisei.name
		}
	}
	return name, att
}
func maxDef(b *Busho) (name string, def int, cost int) {
	for _, tekisei := range b.tekisei {
		if tekisei.def > def {
			def = tekisei.def
			name = tekisei.name
			cost = tekisei.defCost
		}
	}
	return
}

func shubetsu(s *goquery.Selection) string {
	kind := []string{"0", "将", "剣", "忍", "文", "姫", "天", "童"}
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

//Hei 長槍,長弓,精鋭騎馬
type Hei struct {
	name        string
	mainTekisei string
	subTekisei  string
	att         int
	def         int
}

type heiTekisei struct {
	name    string
	att     int
	attCost int
	def     int
	defCost int
}

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
	tekisei  []heiTekisei
	//heilist  []hei
	//nagayari  hei
	//nagayumi  hei
	//seieikiba hei
}

//NewBusho ...
func NewBusho(s *goquery.Selection, h []Hei) *Busho {
	var Busho Busho
	Busho.init(s, h)
	return &Busho
}
func sjis2utf8(str string) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}
func utf82sjis(str string) string {
	ret, _ := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewEncoder()))
	return string(ret)
}
func (b *Busho) init(s *goquery.Selection, H []Hei) {
	if len(s.Find(".ig_card_cost").Text()) == 0 {
		b.cost = 193
	} else {
		b.cost, _ = strconv.ParseFloat(s.Find(".ig_card_cost").Text(), 64)
	}
	b.no = s.Find(".ig_card_cardno").Text()
	//b.name, _ = sjis2utf8(s.Find(".ig_card_name").Text())
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
	for i := 0; i < len(H); i++ {
		var h heiTekisei
		h.set(b, &H[i], b.yari.hosei, b.kiba.hosei, b.yumi.hosei, b.heiki.hosei)
		b.tekisei = append(b.tekisei, h)
	}
}

func (h *heiTekisei) set(b *Busho, H *Hei, hoseiYari int, hoseiKiba int, hoseiYumi int, hoseiHeiki int) {
	cost := b.cost
	if cost == 0 {
		cost = 0.1
	}
	h.name = H.name
	var hosei int
	switch H.mainTekisei {
	case "馬":
		hosei = hoseiKiba
	case "槍":
		hosei = hoseiYari
	case "弓":
		hosei = hoseiYumi
	default:
		hosei = hoseiHeiki
	}
	switch H.subTekisei {
	case "馬":
		hosei = (hosei + hoseiKiba) / 2
	case "槍":
		hosei = (hosei + hoseiYari) / 2
	case "弓":
		hosei = (hosei + hoseiYumi) / 2
	case "器":
		hosei = (hosei + hoseiHeiki) / 2
	}
	h.att = (b.att + b.comno*H.att) * hosei / 100
	h.attCost = int(float64(h.att) / cost)
	h.def = (b.def + b.comno*H.def) * hosei / 100
	h.defCost = int(float64(h.def) / cost)
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
