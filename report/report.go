package report

import (
	"github.com/xuri/excelize/v2"
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"gorm.io/gorm"
	"log"
	"strings"
)

// ranking
const RankingSQL = `select
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

type Ranking struct {
	Count int `gorm:"column:c"`
	Name  string
}

// report
const ReportSQL = `select
    l.c,
    t2.ticker,
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
    t.ticker asc,
    t.tweet_at asc
`

type Report struct {
	Count     int `gorm:"column:c"`
	Ticker    string
	Fund      string
	Name      string
	TwitterId string `gorm:"column:twitter_id"`
	Comment   string
	TweetAt   string `gorm:"column:tweet_at"`
}

func read(ds string) ([]Ranking, []Report, error) {
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

	var rankings []Ranking
	db.Raw(RankingSQL).Scan(&rankings)

	// reports
	var reports []Report
	db.Raw(ReportSQL).Scan(&reports)

	return rankings, reports, nil
}

// output filename
const ReportFile = "./foy2022-tw.xlsx"
const OldName = "Sheet1"
const RankingSheet = "順位"
const CountLabel = "票数"
const FundLabel = "ファンド名"
const NameLabel = "名前"
const IdLabel = "ID"
const CommentLabel = "コメント"
const TimestampLabel = "時刻"
const Sep = "%0D%0A"

func col(col int, row int) string {
	name, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		log.Fatal(err)
	}
	return name
}

func lines(str string) (string, int) {
	if strings.Index(str, Sep) == -1 {
		return str, 1
	}
	lines := strings.Split(str, Sep)
	n := len(lines)
	return strings.Join(lines, "\n"), n
}

func width(f *excelize.File, sheet string, col string, x float64) {
	w := 9.14
	f.SetColWidth(sheet, col, col, w*x)
}

func out(rankings []Ranking, reports []Report) error {
	// open Excel file
	f := excelize.NewFile()
	f.SetSheetName(OldName, RankingSheet)

	// box style
	style, err := f.NewStyle(&excelize.Style{Border: []excelize.Border{
		{Type: "top", Color: "000000", Style: 1},
		{Type: "left", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}})
	if err != nil {
		log.Fatal(err)
	}

	// ranking
	f.SetCellStr(RankingSheet, "A1", CountLabel)
	f.SetCellStr(RankingSheet, "B1", FundLabel)
	width(f, RankingSheet, "B", 10.0)
	for i, r := range rankings {
		f.SetCellValue(RankingSheet, col(1, i+2), r.Count)
		f.SetCellValue(RankingSheet, col(2, i+2), r.Name)
		f.SetCellStyle(RankingSheet, col(1, 1), col(2, i+2), style)
	}

	prev := ""
	n := 0
	for _, r := range reports {
		ticker := r.Ticker
		if ticker != prev {
			n = 0
			f.NewSheet(ticker)
			f.SetCellValue(ticker, "A1", r.Count)
			f.SetCellValue(ticker, "B1", "票")
			f.SetCellValue(ticker, "C1", r.Fund)
			f.SetCellValue(ticker, "A2", NameLabel)
			width(f, ticker, "A", 3.0)
			f.SetCellValue(ticker, "B2", IdLabel)
			width(f, ticker, "B", 2.5)
			f.SetCellValue(ticker, "C2", CommentLabel)
			width(f, ticker, "C", 12.0)
			f.SetCellValue(ticker, "D2", TimestampLabel)
			width(f, ticker, "D", 2.5)
			prev = ticker
		}
		f.SetCellValue(ticker, col(1, n+3), r.Name)
		f.SetCellValue(ticker, col(2, n+3), r.TwitterId)
		com, x := lines(r.Comment)
		f.SetCellValue(ticker, col(3, n+3), com)
		if x > 1 {
			h, err := f.GetRowHeight(ticker, n+3)
			if err != nil {
				log.Fatal(err)
			}
			f.SetRowHeight(ticker, n+3, h*float64(x))
		}
		f.SetCellValue(ticker, col(4, n+3), r.TweetAt)
		f.SetCellStyle(ticker, col(1, 2), col(4, n+3), style)
		n = n + 1
	}

	// close Excel file
	if err := f.SaveAs(ReportFile); err != nil {
		log.Fatal(err)
	}

	return nil
}

func Output(ds string) {
	// read dababase
	rankings, reports, err := read(ds)
	if err != nil {
		log.Fatal(err)
	}

	// output excel
	err = out(rankings, reports)
	if err != nil {
		log.Fatal(err)
	}
}
