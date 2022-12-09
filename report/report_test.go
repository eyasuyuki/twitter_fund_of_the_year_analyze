package report

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestRead(t *testing.T) {
	lankings, reports, err := read("../foy2022.db")
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
