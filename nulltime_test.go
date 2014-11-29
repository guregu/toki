package toki_test

import (
	"encoding/json"
	"testing"

	"github.com/guregu/toki"
)

var jsonTable = map[string]string{
	`"12:34"`: `"12:34"`,
	`""`:      "null",
}

func TestMarshalNullJSON(t *testing.T) {
	for given, expected := range jsonTable {
		time := toki.NullTime{}
		time.UnmarshalJSON([]byte(given))
		data, err := json.Marshal(time)
		if err != nil {
			t.Errorf("MarshalJSON: error: %s", err)
		}
		result := string(data)
		if result != expected {
			t.Errorf("MarshalJSON: %s ≠ %s (%#v)", result, expected, time)
		}
	}
}

func TestMarshalText(t *testing.T) {
	time := toki.MustParseNullTime("")
	if len(time.String()) != 0 {
		t.Errorf("MarshalText: unexpected value %s", time.String())
	}
}

func TestMustParseNullTime(t *testing.T) {
	panicked := false
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
		if !panicked {
			t.Error("MustParseNullTime: didn't panic")
		}
	}()

	for given, expected := range parseTable {
		time := toki.MustParseNullTime(given)
		result := time.String()
		if result != expected {
			t.Errorf("MustParseNullTime: %s → expected %s, got: %s", given, expected, result)
		}
	}

	time := toki.MustParseNullTime("")
	if time.Valid {
		t.Error("MustParseNullTime: empty string isn't null")
	}

	toki.MustParseNullTime("invalid input")
}

func TestNullScan(t *testing.T) {
	expected := toki.MustParseNullTime("12:34")
	byteTime := toki.NullTime{}
	err := byteTime.Scan([]byte("12:34:00"))
	if err != nil {
		t.Errorf("Scan: error: %v", err)
	}
	if byteTime != expected {
		t.Errorf("Scan: %#v ≠ %#v", byteTime, expected)
	}

	nullTime := toki.NullTime{}
	err = nullTime.Scan(nil)
	if err != nil {
		t.Errorf("Scan: error: %v", err)
	}
	if nullTime.Valid {
		t.Errorf("Scan: not null, got %s", nullTime.String())
	}

	emptyTime := toki.NullTime{}
	err = emptyTime.Scan("")
	if err != nil {
		t.Errorf("Scan: error: %v", err)
	}
	if emptyTime.Valid {
		t.Errorf("Scan: not null, got %s", emptyTime.String())
	}

	emptyByteTime := toki.NullTime{}
	err = emptyByteTime.Scan([]byte{})
	if err != nil {
		t.Errorf("Scan: error: %v", err)
	}
	if emptyByteTime.Valid {
		t.Errorf("Scan: not null, got %s", emptyByteTime.String())
	}
}

func TestNullValue(t *testing.T) {
	time := toki.MustParseNullTime("12:34")
	v, err := time.Value()
	if string(v.([]byte)) != "12:34" {
		t.Errorf("Value: %v ≠ 12:34", v)
	}
	if err != nil {
		t.Errorf("Value: error: %v", err)
	}

	nullTime := toki.MustParseNullTime("")
	v, err = nullTime.Value()
	if v != nil {
		t.Errorf("Value: %v ≠ 12:34", v)
	}
	if err != nil {
		t.Errorf("Value: error: %v", err)
	}
}

func BenchmarkNull(b *testing.B) {
	for n := 0; n < b.N; n++ {
		t := toki.NullTime{}
		t.UnmarshalJSON([]byte{'n', 'u', 'l', 'l'})
	}
}

func BenchmarkNormalTime(b *testing.B) {
	for n := 0; n < b.N; n++ {
		t := toki.NullTime{}
		t.UnmarshalJSON([]byte{'1', '2', ':', '3', '0'})
	}
}

func BenchmarkEmptyTime(b *testing.B) {
	for n := 0; n < b.N; n++ {
		t := toki.NullTime{}
		t.UnmarshalJSON([]byte{'"', '"'})
	}
}

func BenchmarkMarshal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		t := toki.NullTime{toki.Time{12, 34, 00}, true}
		t.MarshalJSON()
	}
}

func BenchmarkMarshalNull(b *testing.B) {
	for n := 0; n < b.N; n++ {
		t := toki.NullTime{toki.Time{}, false}
		t.MarshalJSON()
	}
}
