package toki

import (
	"database/sql/driver"
)

type NullTime struct {
	Time
	Valid bool
}

func (t *NullTime) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	return t.Time.UnmarshalText(text)
}

func (t *NullTime) Scan(src interface{}) error {
	switch x := src.(type) {
	case []byte:
		if len(x) == 0 {
			t.Valid = false
			return nil
		}
	case string:
		if x == "" || str == "null" {
			t.Valid = false
			return nil
		}
	case nil:
		t.Valid = false
		return nil
	}
	return t.Time.Scan(src)
}

func (t NullTime) MarshalText() (text []byte, err error) {
	if !t.Valid {
		return []byte{}, nil
	}
	return t.Time.MarshalText()
}

func (t NullTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time.MarshalText()
}

func ParseNullTime(text string) (NullTime, error) {
	t := &NullTime{}
	err := t.UnmarshalText([]byte(text))
	return *t, err
}

func MustParseNullTime(text string) NullTime {
	t, err := ParseNullTime(text)
	if err != nil {
		panic(err)
	}
	return t
}
