package toki

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Time struct {
	Hours   int
	Minutes int
	Seconds int
	// TODO: millis?
}

func (t *Time) UnmarshalText(text []byte) error {
	parts := strings.Split(string(text), ":")
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return err
		}
		switch i {
		case 0:
			t.Hours = n
		case 1:
			t.Minutes = n
		case 2:
			t.Seconds = n
		}
	}
	return nil
}

func (t *Time) Scan(src interface{}) error {
	switch x := src.(type) {
	case []byte:
		return t.UnmarshalText(x)
	case string:
		return t.UnmarshalText([]byte(x))
	case time.Time:
		t.Hours = x.Hour()
		t.Minutes = x.Minute()
		t.Seconds = x.Second()
		return nil
	}
	return fmt.Errorf("unsupported type: %T", src)
}

func (t Time) MarshalText() (text []byte, err error) {
	if t.Seconds == 0 {
		return []byte(fmt.Sprintf("%02d:%02d", t.Hours, t.Minutes)), nil
	}
	return []byte(fmt.Sprintf("%02d:%02d:%02d", t.Hours, t.Minutes, t.Seconds)), nil
}

func (t Time) Value() (driver.Value, error) {
	return t.MarshalText()
}

func (t Time) String() string {
	text, _ := t.MarshalText()
	return string(text)
}

func ParseTime(text string) (Time, error) {
	t := &Time{}
	err := t.UnmarshalText([]byte(text))
	return *t, err
}

func MustParseTime(text string) Time {
	t, err := ParseTime(text)
	if err != nil {
		panic(err)
	}
	return t
}
