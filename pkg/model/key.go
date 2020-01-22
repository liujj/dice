package model

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"dice/pkg/utils"
	"time"
)

type Key struct {
	ID         uint         `json:"-" gorm:"primary_key"`
	Ak         string       `json:"ak"`
	Sk         string       `json:"sk"`
	Balance    float64      `json:"balance"`
	Peak       float64      `json:"peak"`
	TxCount    uint         `json:"tx_count"`
	Remark     string       `json:"remark"`
	Vaild      bool         `json:"vaild"`
	BankruptAt *JsonTime    `json:"bankrupt_at"`
	CreatedAt  *JsonTime    `json:"create_at"`
	mu         sync.RWMutex `json:"-" gorm:"-"`
}

func (db *Key) TableName() string {
	return "key"
}

var KC = struct {
	Map   map[string]*Key
	Mutex sync.RWMutex
}{
	Map: make(map[string]*Key),
}

func NewKey() (*Key, error) {
	newKey := &Key{
		Ak:         utils.GenerateRamdom(8),
		Sk:         utils.GenerateSecret(),
		Balance:    20,
		Vaild:      true,
		BankruptAt: nil,
	}
	dbResult := Gdb.Create(newKey)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	} else if dbResult.RowsAffected == 1 {
		return newKey, nil
	} else {
		return newKey, errors.New("failed")
	}
}

func QueryKey(ak string) (*Key, error) {
	k := new(Key)
	dbResult := Gdb.Where("ak = ?", ak).First(k)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	} else if k.Ak == "" {
		return nil, errors.New(fmt.Sprintf("Access Key %s Not Found", ak))
	} else {
		return k, nil
	}
}

func ListKeys(fr, to int) (ks []*Key, e error) {
	dbErr := Gdb.Offset(fr).Limit(to - fr).Order("`id` DESC").Find(&ks)
	return ks, dbErr.Error
}

func SetKey(ak string, param map[string]string) error {
	key := &Key{}
	updateMap := make(map[string]interface{})
	if k, e := param["vaild"]; e && (k == "true" || k == "false") {
		updateMap["vaild"] = k == "true"
	}
	if k, e := param["remark"]; e {
		updateMap["remark"] = k
	}
	if k, e := param["bankruptAt"]; e {
		if k == "null" {
			updateMap["bankrupt_at"] = nil
		} else {
			tempTime, err := time.ParseInLocation(TimeFormart, k, time.Local)
			if err != nil {
				return errors.New("bankruptAt format: 2006-01-02 15:04:05")
			}
			updateMap["bankrupt_at"] = &JsonTime{tempTime}
		}
	} else if k, e := param["bankrupt"]; e {
		if k == "-1" {
			updateMap["bankrupt_at"] = nil
		} else {
			tempInt, err := strconv.Atoi(k)
			if err != nil {
				return errors.New("bankrupt shall be a number")
			}
			if tempInt < 0 {
				return errors.New("bankrupt shall be a positive number or -1")
			}
			updateMap["bankrupt_at"] = &JsonTime{time.Now().Add(time.Duration(tempInt) * 1e9)}
		}
	}
	//使用对象更新时gorm无法更新到零值，故必须使用map强制
	dbResult := Gdb.Model(key).Where("ak = ?", ak).Updates(updateMap)
	KC.Mutex.Lock()
	delete(KC.Map, ak)
	KC.Mutex.Unlock()
	if dbResult.Error != nil {
		return dbResult.Error
	} else if dbResult.RowsAffected == 0 {
		return errors.New(fmt.Sprintf("Modify Failed at Access Key %s", ak))
	} else {
		return nil
	}
}

func DelKey(ak string) error {
	k, _ := QueryKey(ak)
	if k == nil || k.ID == 0 {
		return errors.New(fmt.Sprintf("Access Key %s is not Found", ak))
	}
	dbResult := Gdb.Delete(k)
	KC.Mutex.Lock()
	delete(KC.Map, ak)
	KC.Mutex.Unlock()
	if dbResult.Error != nil {
		return dbResult.Error
	} else if dbResult.RowsAffected == 0 {
		return errors.New(fmt.Sprintf("Access Key %s is not Deleted", ak))
	} else {
		return nil
	}
}

func FetchKey(ak string) (*Key, error) {
	var k *Key
	var ok bool
	KC.Mutex.RLock()
	k, ok = KC.Map[ak]
	KC.Mutex.RUnlock()
	//缓存不存在则加缓存
	if !ok {
		k, _ = QueryKey(ak)
		KC.Mutex.Lock()
		KC.Map[ak] = k
		KC.Mutex.Unlock()
	}

	if k == nil {
		//实际为key不存在
		return nil, errors.New("key does not exist")
	}
	return k, nil
}

func ClearKeyCache() {
	KC.Mutex.Lock()
	KC.Map = make(map[string]*Key)
	KC.Mutex.Unlock()
}

func (k *Key) GetHistory(fromLast, step int) []Roll {
	var rs []Roll
	Gdb.Where("key = ?", k.Ak).Limit(step).Offset(fromLast).Order("id DESC").Find(&rs)
	return rs
}
