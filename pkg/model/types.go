package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	TimeFormartZone = "2006-01-02 15:04:05 -0700"
	TimeFormart     = "2006-01-02 15:04:05"
	TimeFormartNums = "060102150405"
	TimeFormartDate = "2006-01-02"
)

type JsonTime struct {
	time.Time
}

func (t JsonTime) MarshalJSON() ([]byte, error) {
	if t.Unix() == 0 {
		return []byte("null"), nil
	}
	b := make([]byte, 0, len(TimeFormartZone)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, TimeFormartZone)
	b = append(b, '"')
	return b, nil
}

func (ts JsonTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = ts.Time
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti.Format(TimeFormartZone), nil
}

func (ts *JsonTime) Scan(v interface{}) error {
	value_t, ok := v.(time.Time)
	if ok {
		*ts = JsonTime{value_t}
		return nil
	}
	value_byte, ok := v.([]byte)
	if ok {
		t, err := time.Parse(TimeFormartZone, string(value_byte))
		if err != nil {
			return fmt.Errorf("can not convert %v to timestamp", v)
		}
		*ts = JsonTime{t}
		return nil
	}
	value_str, ok := v.(string)
	if ok {
		t, err := time.Parse(TimeFormartZone, value_str)
		if err != nil {
			return fmt.Errorf("can not convert %v to timestamp", v)
		}
		*ts = JsonTime{t}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
