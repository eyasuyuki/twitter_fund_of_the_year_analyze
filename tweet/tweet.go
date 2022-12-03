package tweet

type Tweet struct {
	Id        string
	Name      string
	Fund      string
	Comment   string
	Timestamp string
}

type Fund struct {
	Ticker string
	Tweets []*Tweet
	Count  int64
}
