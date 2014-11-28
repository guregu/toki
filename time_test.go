package toki_test

import (
	"testing"
	"time"

	"github.com/guregu/toki"
)

var parseTable = map[string]string{
	"12:34:56": "12:34:56",
	"12:34":    "12:34",
	"12":       "12:00",
	"1":        "1:00",
}

func TestParseTime(t *testing.T) {
	for given, expected := range parseTable {
		time, err := toki.ParseTime(given)
		if err != nil {
			t.Errorf("ParseTime: %s → expected %s, got error: %v", given, expected, time)
		}
		result := time.String()
		if result != expected {
			t.Errorf("ParseTime: %s → expected %s, got: %s", given, expected, result)
		}
	}

	_, err := toki.ParseTime("12:abcdef")
	if err == nil {
		t.Errorf("ParseTime: expected error, got: nil")
	}
}

func TestMustParseTime(t *testing.T) {
	panicked := false
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
		if !panicked {
			t.Error("MustParseTime: didn't panic")
		}
	}()

	for given, expected := range parseTable {
		time := toki.MustParseTime(given)
		result := time.String()
		if result != expected {
			t.Errorf("MustParseTime: %s → expected %s, got: %s", given, expected, result)
		}
	}

	toki.MustParseTime("invalid input")
}

func TestScan(t *testing.T) {
	expected := toki.MustParseTime("12:34")

	strTime := toki.Time{}
	err := strTime.Scan("12:34")
	if err != nil {
		t.Errorf("Scan: error: %v", err)
	}
	if strTime != expected {
		t.Errorf("Scan: %#v ≠ %#v", strTime, expected)
	}

	byteTime := toki.Time{}
	err = byteTime.Scan([]byte("12:34:00"))
	if err != nil {
		t.Errorf("Scan: error: %v", err)
	}
	if byteTime != expected {
		t.Errorf("Scan: %#v ≠ %#v", byteTime, expected)
	}

	timeTime := toki.Time{}
	err = timeTime.Scan(time.Date(1992, 2, 23, 12, 34, 0, 0, time.UTC))
	if err != nil {
		t.Errorf("Scan: error: %v", err)
	}
	if timeTime != expected {
		t.Errorf("Scan: %#v ≠ %#v", timeTime, expected)
	}

	errTime := toki.Time{}
	err = errTime.Scan(42)
	if err == nil {
		t.Errorf("Scan: expected error, got: nil")
	}
}

func TestValue(t *testing.T) {
	time := toki.MustParseTime("12:34")
	v, err := time.Value()
	if string(v.([]byte)) != "12:34" {
		t.Errorf("Value: %v ≠ 12:34", v)
	}
	if err != nil {
		t.Errorf("Value: error: %v", err)
	}
}
