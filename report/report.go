package report

import (
	"github.com/xuri/excelize/v2"
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"gorm.io/gorm"
	"log"
)

// lanking
const LankingSQL = `select
    count(tweets.ticker) c,
    tickers.name
from
    tweets
        inner join tickers
                   on tweets.ticker=tickers.ticker
group by
    tweets.ticker
order by
    c desc
`

type Lanking struct {
	Count int64 `gorm:"column:c"`
	Name  string
}

// report
const ReportSQL = `select
    l.c, 
    t2.name fund,
    t.name,
    t.twitter_id,
    t.comment,
    t.tweet_at
from
    tweets t
    inner join tickers t2
        on t.ticker = t2.ticker
    inner join (select
                    ticker,
                    count(ticker) c
                from
                    tweets
                group by
                    ticker) l
        on t.ticker = l.ticker
order by
    l.c desc,
    t.tweet_at asc
`

type Report struct {
	Count     int64 `gorm:"column:c"`
	Fund      string
	Name      string
	TwitterId string `gorm:"column:twitter_id"`
	Comment   string
	TweetAt   string `gorm:"column:tweet_at"`
}

func read(ds string) ([]Lanking, []Report, error) {
	// open sqlite
	db, err := gorm.Open(sqlite.Open(ds), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	d, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	defer d.Close()

	var lankings []Lanking
	db.Raw(LankingSQL).Scan(&lankings)

	// reports
	var reports []Report
	db.Raw(ReportSQL).Scan(&reports)

	return lankings, reports, nil
}

// output filename
const ReportFile = "../foy2022-tw.xlsx"

func out(lankings []Lanking, reports []Report) error {
	// open Excel file
	f := excelize.NewFile()
	// close Excel file
	if err := f.SaveAs(ReportFile); err != nil {
		log.Fatal(err)
	}

	return nil
}

func Output(ds string) {
	// read dababase
	lankings, reports, err := read(ds)
	if err != nil {
		log.Fatal(err)
	}

	// output excel
	err = out(lankings, reports)
	if err != nil {
		log.Fatal(err)
	}
}
