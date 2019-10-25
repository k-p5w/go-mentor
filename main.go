package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

// Resultitem is 結果表示用
type Resultitem struct {
	Title string
	Data  []Clippage
}

// AllItem is 全件の格納用
type AllItem struct {
	AllData []Clippage
}

// Clippage is ページの情報
type Clippage struct {
	Title  string
	Urltxt string
	Rank   string
	Point  int
}

// Allpageinfo is 全データ
var Allpageinfo []Clippage

// main is メイン表示用の関数
func main() {

	url := "https://www.gov-online.go.jp/useful/index.html"
	// htmlを読み込んで一覧を作る
	Allpageinfo = readhtmlpage(url)

	fs := http.FileServer(http.Dir("./tmp"))

	http.HandleFunc("/get", viewHandler)

	http.Handle("/data/", http.StripPrefix("/data/", fs))
	http.HandleFunc("/", initHandler)

	// HEROKUで動かすためにポートは取り出すようにする
	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}
	http.ListenAndServe(":"+port, nil)

}

// initHandler is 初期表示
func initHandler(w http.ResponseWriter, r *http.Request) {

	page := Resultitem{"～初期表示～", Allpageinfo}
	//HTMLを読み込む
	tmpl, err := template.ParseFiles("./tmp/allview.html") // ParseFilesを使う
	if err != nil {
		panic(err)
	}
	//描画させる
	err = tmpl.Execute(w, page)
	if err != nil {
		panic(err)
	}
}

// viewHandler is postでランダム表示を実現するため
func viewHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("viewHandler-go!")
	// 表の表示件数
	maxitem := 10

	// getパラメータの解析
	q := r.URL.Query()
	getval := q.Get("item")
	if len(getval) > 0 {
		newval, _ := strconv.Atoi(getval)
		// 最大件数未満なら書き換える
		if newval <= maxitem {
			maxitem = newval
		}

	}

	choisdata := make([]Clippage, maxitem)
	selected := make([]int, maxitem)
	selitem := 0
	for choisidx := 0; choisidx < 100; choisidx++ {

		// ランダム値を取得する
		var inf64 int64
		inf64 = int64(choisidx + 1)
		r := rand.New(rand.NewSource(67))
		r.Seed(time.Now().UnixNano() + inf64)
		idx := r.Intn(len(Allpageinfo))

		IsItemExists := func(targetidx int) bool {
			for _, v := range selected {
				// 同じ値が取得された場合は取り直す
				if targetidx == v {
					return true
				}
			}
			return false
		}
		// 見つかるたびにポイントを加算していく
		Allpageinfo[idx].Point++
		// 同じ値が存在していればやり直し
		if IsItemExists(idx) {
			continue
		}
		// 配列にアイテムをセットする
		choisdata[selitem].Title = Allpageinfo[idx].Title
		choisdata[selitem].Urltxt = Allpageinfo[idx].Urltxt
		choisdata[selitem].Rank = getRank(idx)

		selected[selitem] = idx
		selitem++

		// 選択が終了した場合
		if selitem == len(choisdata) {
			break
		}
	}

	page := Resultitem{"ランダム表示", choisdata}
	//HTMLを読み込む
	tmpl, err := template.ParseFiles("./tmp/simple.html") // ParseFilesを使う
	if err != nil {
		panic(err)
	}
	//描画させる
	err = tmpl.Execute(w, page)
	if err != nil {
		panic(err)
	}
}

// getRank is ランクを返す
func getRank(index int) string {
	// 表示した回数をもとにランキング化して色を設定する
	ret := "flex-item card-common"
	rank3 := "flex-item card-bronze"
	rank2 := "flex-item card-silver"
	rank1 := "flex-item card-gold"
	rank0 := "flex-item card-secret"
	max := len(Allpageinfo) - 1
	// コピーして並び替える
	sortdata := make([]Clippage, len(Allpageinfo))
	cLen := copy(sortdata, Allpageinfo)
	fmt.Println(cLen) // 5

	sort.Slice(sortdata, func(i, j int) bool {
		return sortdata[i].Point > sortdata[j].Point
	})

	// 1-3位のポイント以上ある？
	if sortdata[0].Point <= Allpageinfo[index].Point {
		ret = rank1
	} else {
		if sortdata[1].Point <= Allpageinfo[index].Point {
			ret = rank2
		} else {
			if sortdata[2].Point <= Allpageinfo[index].Point {
				ret = rank3
			} else {
				// 最下位の場合
				if sortdata[max].Point == (Allpageinfo[index].Point-1) {
					ret = rank0
				}
			}
		}

	}

	return ret
}
