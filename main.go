package main

import (
	"database/sql"
	"github.com/PuerkitoBio/goquery"
	"github.com/eyasuyuki/twitter_fund_of_the_year_analyze/config"
	"github.com/eyasuyuki/twitter_fund_of_the_year_analyze/report"
	"github.com/eyasuyuki/twitter_fund_of_the_year_analyze/tweet"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"strings"
)

// extract fund and comment
func fundComment(cf *config.Config, tweet string) (string, string) {
	var fund, comment string
	if tweet == "" {
		return "", ""
	}

	lines := strings.Split(tweet, "\n")
	commentFound := false
	for i, line := range lines {
		if line == "" {
			continue
		}
		line = strings.TrimSpace(line)
		if fund == "" && strings.Index(line, "🥇") > -1 {
			fund = lines[i+1]
		}
		if strings.Index(line, "→") > -1 {
			commentFound = true
			ls := strings.Split(line, "→")
			if len(ls) > 1 && ls[1] != "" && strings.Index(line, cf.PageUrl) == -1 {
				comment = ls[1]
			}
		} else if commentFound && line != "" && strings.Index(line, cf.PageUrl) == -1 {
			if comment != "" {
				comment = comment + "%0D%0A"
			}
			comment = comment + line
		}
	}

	return fund, comment
}

const PageTweet = 25

// tweets table
const CreateTweets = `create table tweets (id integer not null primary key, ticker text, twitter_id text, name text, fund text, comment text, tweet_at text)`
const InsertTweets = `insert into tweets (id, ticker, twitter_id, name, fund, comment, tweet_at) values (?, ?, ?, ?, ?, ?, ?)`

// ticker table
const CreateTickers = `create table tickers (ticker not null primary key, name text)`
const SelectTickers = `select * from tickers where ticker = ?`
const InsertTickers = `insert into tickers (ticker, name) values (?, ?)`

func main() {
	cf := config.NewConfig("")
	urls := []string{cf.TogetterUrl,
		cf.TogetterUrl + "?page=2",
		cf.TogetterUrl + "?page=3",
		cf.TogetterUrl + "?page=4",
	}
	n := 0

	// delete db
	os.Remove(cf.DatabaseName)
	// open
	db, err := sql.Open("sqlite3", cf.DatabaseName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// create table
	_, err = db.Exec(CreateTweets)
	if err != nil {
		log.Printf("%q: %s\n", err, CreateTweets)
		return
	}
	_, err = db.Exec(CreateTickers)
	if err != nil {
		log.Printf("%q: %s\n", err, CreateTweets)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	s1, err := tx.Prepare(InsertTweets)
	if err != nil {
		log.Fatal(err)
	}
	s2, err := tx.Prepare(InsertTickers)
	if err != nil {
		log.Fatal(err)
	}

	tweets := []tweet.Tweet{}
	for _, url := range urls {
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			log.Fatalf("status code err: %d %s", res.StatusCode, res.Status)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		// fundComment fund name, comment
		doc.Find(".tweet").Each(func(i int, s *goquery.Selection) {
			fund, comment := fundComment(cf, s.Text())
			t := tweet.Tweet{Fund: fund, Comment: comment}
			tweets = append(tweets, t)
		})

		doc.Find(".user_link").Each(func(i int, s *goquery.Selection) {
			name := s.Find("strong").Text()
			id := s.Find(".status_name").Text()
			if name != "" && id != "" {
				tweets[n*PageTweet+i].Name = name
				tweets[n*PageTweet+i].Id = id

			}
		})

		doc.Find(".status").Each(func(i int, s *goquery.Selection) {
			ts := s.Find("a").Text()
			index := n*PageTweet + i
			tweets[index].Timestamp = strings.TrimSpace(ts)
			// insert
			_, err = s1.Exec(index, tweets[index].Ticker(), tweets[index].Id, tweets[index].Name, tweets[index].Fund, tweets[index].Comment, tweets[index].Timestamp)
			if err != nil {
				log.Fatal(err)
			}
			rows, err := db.Query(SelectTickers, tweets[index].Ticker())
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			if !rows.Next() {
				s2.Exec(tweets[index].Ticker(), tweets[index].Fund)
			}
		})

		n = n + 1
	}
	// commit
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	report.Output(cf, cf.DatabaseName) //TEST

	//jsonText, err := json.Marshal(tweets)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("%v", string(jsonText))
}
