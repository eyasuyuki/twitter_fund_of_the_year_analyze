package tweet

import "strings"

type Tweet struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Fund      string `json:"fund"`
	Comment   string `json:"comment"`
	Timestamp string `json:"timestamp"`
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
