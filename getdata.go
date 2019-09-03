package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// readhtmlpage is htmlページのスクレイピング
func readhtmlpage(url string) []Clippage {

	var tmpPage Clippage
	allpage := make([]Clippage, 0)

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

	doc.Find("#article > div.body > div").Each(func(index int, s *goquery.Selection) {
		s.Find("li").Each(func(i int, sec *goquery.Selection) {
			title := sec.Find("a").Text()
			urllink, _ := sec.Find("a").Attr("href")
			itemkeyword := sec.Find("div.keyword").Text()
			fmt.Printf("Review %d: %s - %s\n", i, urllink, itemkeyword)
			tmpPage.Title = strings.TrimSpace(title)
			tmpPage.Urltxt = urllink
			allpage = append(allpage, tmpPage)
		})
	})

	return allpage
}
