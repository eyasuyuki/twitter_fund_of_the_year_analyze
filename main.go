package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

func analyze(tweet string) (string, string) {
	var fund, comment string
	if tweet == "" {
		return "", ""
	}
	lines := strings.Split(tweet, "\n")
	for i, line := range lines {
		if fund == "" && strings.Index(line, "🥇") > -1 {
			fund = lines[i+1]
		}
		if comment == "" && strings.Index(line, "→") > -1 {

			break
		}
	}
	return fund, comment
}

func main() {
	urls := []string{"https://togetter.com/li/1980458",
		"https://togetter.com/li/1980458?page=2",
		"https://togetter.com/li/1980458?page=3",
		"https://togetter.com/li/1980458?page=4",
	}
	for _, url := range urls {
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			log.Fatal("status code err: %d %s", res.StatusCode, res.Status)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		doc.Find(".user_link").Each(func(i int, s *goquery.Selection) {
			name := s.Find("strong").Text()
			id := s.Find(".status_name").Text()
			if name != "" && id != "" {
				fmt.Printf("%d: %s %s\n", i, name, id)
			}
		})
		doc.Find(".tweet").Each(func(i int, s *goquery.Selection) {
			fund, comment := analyze(s.Text())
			fmt.Printf("%d: %s %s\n", i, fund, comment)
		})
		doc.Find(".status").Each(func(i int, s *goquery.Selection) {
			ts := s.Find("a").Text()
			fmt.Printf("%d: %s\n", i, strings.TrimSpace(ts))
		})
	}
}
