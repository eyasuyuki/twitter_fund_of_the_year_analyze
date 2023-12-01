package config

import "testing"

func testConfig(t *testing.T) {
	c := NewConfig("")
	if c == nil {
		t.Error("Config is nil")
	}
	// Year
	if c.Year != "2023" {
		t.Errorf("Year expect 2023, but %v", c.Year)
	}
	// PageUrl
	if c.PageUrl != "https://www.fundoftheyear.jp/2023/tweet.html" {
		t.Errorf("PageUrl expect https://www.fundoftheyear.jp/2023/tweet.html, but %v", c.PageUrl)
	}
	// ToggerUrl
	if c.TogetterUrl != "https://togetter.com/li/2268120" {
		t.Errorf("ToggerUrl expect https://togetter.com/li/2268120, but %v", c.TogetterUrl)
	}
	// DatabaseName
	if c.DatabaseName != "foy2023.db" {
		t.Errorf("DatabaseName exepct foy2023.db, but %v", c.DatabaseName)
	}
	// ReportFile

}
