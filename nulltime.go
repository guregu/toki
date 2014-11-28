package toki

import (
	"database/sql/driver"
)

type NullTime struct {
	Time
	Valid bool
}

func (t *NullTime) UnmarshalText(text []byte) error {
	t.Valid = true
	str := string(text)
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	return t.Time.UnmarshalText(text)
}

func (t *NullTime) UnmarshalJSON(data []byte) error {
	t.Valid = true
	text := string(data)
	if text == `""` || text == "null" {
		t.Valid = false
		return nil
	}
	return t.UnmarshalText(data[1 : len(data)-1])
}

func (t *NullTime) Scan(src interface{}) error {
	t.Valid = true
	switch x := src.(type) {
	case []byte:
		if len(x) == 0 {
			t.Valid = false
			return nil
		}
	case string:
		if x == "" || x == "null" {
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

func (t NullTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	text, _ := t.MarshalText()
	// what is the best way to do this?
	out := make([]byte, 0, len(text)+2)
	out = append(out, '"')
	out = append(out, text...)
	out = append(out, '"')
	return out, nil
}

func (t NullTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time.MarshalText()
}

func (t NullTime) String() string {
	text, _ := t.MarshalText()
	return string(text)
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
