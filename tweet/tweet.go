package tweet

import "strings"

type Tweet struct {
	Id        string
	Name      string
	Fund      string
	Comment   string
	Timestamp string
}

func (t *Tweet) Ticker() string {
	fs := strings.Split(t.Fund, " ")
	if len(fs) > 1 {
		return fs[0]
	} else if len(t.Fund) > 0 {
		return t.Fund
	} else {
		return ""
	}
}
