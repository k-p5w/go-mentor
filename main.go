package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Resultitem is 結果表示用
type Resultitem struct {
	Data []Clippage
}

// AllItem is 全件の格納用
type AllItem struct {
	AllData []Clippage
}

// Clippage is ページの情報
type Clippage struct {
	Title  string
	Urltxt string
}

// Allpageinfo is 全データ
var Allpageinfo []Clippage

// main is メイン表示用の関数
func main() {

	url := "https://www.gov-online.go.jp/useful/index.html"
	// htmlを読み込んで一覧を作る
	Allpageinfo = readhtmlpage(url)

	fs := http.FileServer(http.Dir("./tmp"))

	http.HandleFunc("/get", viewHandler) // mentor初期画面
	http.HandleFunc("/", initHandler)    // mentor初期画面
	http.Handle("/data/", http.StripPrefix("/data/", fs))

	// HEROKUで動かすためにポートは取り出すようにする
	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}
	http.ListenAndServe(":"+port, nil)

}

// initHandler is 初期表示
func initHandler(w http.ResponseWriter, r *http.Request) {

	page := Resultitem{Allpageinfo}
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
		// 同じ値が存在していればやり直し
		if IsItemExists(idx) {
			continue
		}
		// 配列にアイテムをセットする
		choisdata[selitem].Title = Allpageinfo[idx].Title
		choisdata[selitem].Urltxt = Allpageinfo[idx].Urltxt

		selected[selitem] = idx
		selitem++

		// 選択が終了した場合
		if selitem == len(choisdata) {
			break
		}
	}

	page := Resultitem{choisdata}
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
