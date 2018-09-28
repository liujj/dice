package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"tigerMachine/utils"
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
	value_str, ok := v.([]byte)
	if ok {
		t, err := time.Parse(TimeFormartZone, string(value_str))
		if err != nil {
			return fmt.Errorf("can not convert %v to timestamp", v)
		}
		*ts = JsonTime{t}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

type StringSlice []string

func (ss StringSlice) Value() (driver.Value, error) {
	a, err := json.Marshal(ss)
	if err != nil {
		return "", err
	}
	return string(a), nil
}

func (ss *StringSlice) Scan(v interface{}) error {
	value, ok := v.([]byte)
	if ok {
		err := json.Unmarshal(value, ss)
		if err != nil {
			ss = new(StringSlice)
			return err
		}
		return nil
	}
	return fmt.Errorf("can not convert %v to []string", v)
}

type Ip struct {
	uint32
	uint8
}

func NewIp(str string) Ip {
	ip, sub := utils.IpRegion2Long(str)
	ipR := Ip{ip, sub}
	return ipR
}

func (ip Ip) Verify(clientIpStr string) bool {
	if ip.uint8 == 0 {
		return true
	}
	clientIpObj := utils.Ip2Long(clientIpStr)
	if clientIpObj == 0 {
		logrus.Warningf("client ip %s cannot be prased", clientIpStr)
		return false
	}
	mask := 32 - ip.uint8
	if ip.uint32>>mask == clientIpObj>>mask {
		return true
	} else {
		return false
	}
}

func (ip Ip) String() string {
	return utils.Long2IpRegion(ip.uint32, ip.uint8)
}

func (ip Ip) MarshalJSON() ([]byte, error) {
	return []byte("\"" + utils.Long2IpRegion(ip.uint32, ip.uint8) + "\""), nil
}

type IpSlice []Ip

func (ipArr IpSlice) Value() (driver.Value, error) {
	strArr := []string{}
	for _, ip := range ipArr {
		strArr = append(strArr, ip.String())
	}
	a, err := json.Marshal(strArr)
	if err != nil {
		return "", err
	}
	return string(a), nil
}

func (ip *IpSlice) Scan(v interface{}) error {
	value, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("can not convert %v to []string", v)
	}
	strArr := []string{}
	err := json.Unmarshal(value, &strArr)
	if err != nil {
		return err
	}
	var ipTemp Ip
	for _, ipStr := range strArr {
		ipTemp = NewIp(ipStr)
		*ip = append(*ip, ipTemp)
	}
	return nil
}

func (ipArr IpSlice) Verify(clientIp string) error {
	for _, ip := range ipArr {
		if ip.Verify(clientIp) {
			return nil
		}
	}
	return errors.New("ip verify failed")
}
