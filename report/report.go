package report

import (
	"github.com/eyasuyuki/twitter_fund_of_the_year_analyze/config"
	"github.com/xuri/excelize/v2"
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
)

// ranking
const RankingSQL = `WITH CandidateVote AS (
    select count(tweets.ticker) c,
          tickers.name name
   from tweets
            inner join tickers
                       on tweets.ticker = tickers.ticker
   group by tweets.ticker
),
RankedVotes AS (
    SELECT
        name,
        c,
        RANK() OVER (ORDER BY c DESC) AS r
    FROM
        CandidateVote
)
SELECT
    r,
    c,
    name
FROM
    RankedVotes
`

type Ranking struct {
	Rank  int `gorm:"column:r"`
	Count int `gorm:"column:c"`
	Name  string
}

// report
const ReportSQL = `select
	l.r,
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
    inner join (WITH CandidateVote AS (
    select count(tweets.ticker) c,
           tickers.ticker ticker
    from tweets
             inner join tickers
                        on tweets.ticker = tickers.ticker
    group by tweets.ticker
),
     RankedVotes AS (
         SELECT
             ticker,
             c,
             RANK() OVER (ORDER BY c DESC) AS r
         FROM
             CandidateVote
     )
SELECT
    r,
    c,
    ticker
FROM
    RankedVotes) l
        on t.ticker = l.ticker
order by
    l.c desc,
    t.ticker asc,
    t.tweet_at asc
`

type Report struct {
	Rank      int `gorm:"column:r"`
	Count     int `gorm:"column:c"`
	Ticker    string
	Fund      string
	Name      string
	TwitterId string `gorm:"column:twitter_id"`
	Comment   string
	TweetAt   string `gorm:"column:tweet_at"`
}

func read(cf *config.Config, ds string) ([]Ranking, []Report, error) {
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
const OldName = "Sheet1"
const RankingLabel = "順位"
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

func out(cf *config.Config, rankings []Ranking, reports []Report) error {
	// open Excel file
	f := excelize.NewFile()
	f.SetSheetName(OldName, RankingLabel)

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
	f.SetCellStr(RankingLabel, "A1", RankingLabel)
	f.SetCellStr(RankingLabel, "B1", CountLabel)
	f.SetCellStr(RankingLabel, "C1", FundLabel)
	width(f, RankingLabel, "C", 12.5)
	for i, r := range rankings {
		f.SetCellValue(RankingLabel, col(1, i+2), r.Rank)
		f.SetCellValue(RankingLabel, col(2, i+2), r.Count)
		f.SetCellValue(RankingLabel, col(3, i+2), r.Name)
		f.SetCellStyle(RankingLabel, col(1, 1), col(3, i+2), style)
	}

	prev := ""
	n := 0
	offset := 4
	for _, r := range reports {
		rank := strconv.Itoa(r.Rank)
		sheet := rank + "_" + r.Ticker
		if sheet != prev {
			n = 0
			f.NewSheet(sheet)
			f.SetCellValue(sheet, "A1", RankingLabel)
			f.SetCellValue(sheet, "B1", rank)
			f.SetCellValue(sheet, "C1", r.Fund)
			f.SetCellValue(sheet, "A2", CountLabel)
			f.SetCellValue(sheet, "B2", r.Count)
			f.SetCellValue(sheet, "A3", NameLabel)
			width(f, sheet, "A", 2.5)
			f.SetCellValue(sheet, "B3", IdLabel)
			width(f, sheet, "B", 2.5)
			f.SetCellValue(sheet, "C3", CommentLabel)
			width(f, sheet, "C", 12.5)
			f.SetCellValue(sheet, "D3", TimestampLabel)
			width(f, sheet, "D", 3.0)
			prev = sheet
		}
		f.SetCellValue(sheet, col(1, n+offset), r.Name)
		f.SetCellValue(sheet, col(2, n+offset), r.TwitterId)
		com, x := lines(r.Comment)
		f.SetCellValue(sheet, col(3, n+offset), com)
		if x > 1 {
			h, err := f.GetRowHeight(sheet, n+offset)
			if err != nil {
				log.Fatal(err)
			}
			f.SetRowHeight(sheet, n+offset, h*float64(x))
		}
		f.SetCellValue(sheet, col(4, n+offset), r.TweetAt)
		f.SetCellStyle(sheet, col(1, 3), col(4, n+offset), style)
		n = n + 1
	}

	// close Excel file
	if err := f.SaveAs(cf.ReportFile); err != nil {
		log.Fatal(err)
	}

	return nil
}

func Output(cf *config.Config, ds string) {
	// read dababase
	rankings, reports, err := read(cf, ds)
	if err != nil {
		log.Fatal(err)
	}

	// output excel
	err = out(cf, rankings, reports)
	if err != nil {
		log.Fatal(err)
	}
}
