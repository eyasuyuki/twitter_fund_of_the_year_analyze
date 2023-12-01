package report

import (
	"encoding/json"
	"fmt"
	"github.com/eyasuyuki/twitter_fund_of_the_year_analyze/config"
	"log"
	"testing"
)

func TestRead(t *testing.T) {
	cf := config.NewConfig("../config.json")
	lankings, reports, err := read(cf, "../"+cf.DatabaseName)
	if err != nil {
		log.Fatal(err)
	}

	// debug
	js, err := json.Marshal(lankings)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(js))

	js, err = json.Marshal(reports)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(js))
}
