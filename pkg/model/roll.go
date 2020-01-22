package model

import (
	"math/rand"
	"time"
)

type Roll struct {
	ID         uint     `json:"-" gorm:"primary_key"`
	Key        string   `json:"key" gorm:"unique"`
	KeyObj     *Key     `json:"-" gorm:"-"`
	Bet        float64  `json:"bet"`
	RollUnder  uint8    `json:"roll_under"`
	RollRandom uint8    `json:"roll_random"`
	Payout     float64  `json:"payout"`
	Result     bool     `json:"result"`
	Before     float64  `json:"before"`
	After      float64  `json:"after"`
	CreatedAt  JsonTime `json:"roll_time"`
}

func (c *Roll) TableName() string {
	return "roll"
}

func MakeRoll(k *Key, bet float64, under uint8) *Roll {
	if !k.Vaild {
		return nil
	}
	if under < 1 || under > 95 {
		return nil
	}
	if bet > k.Balance || bet < 1 || bet > 5000 {
		return nil
	}
	rand32 := rand.Intn(100)
	rand8 := uint8(rand32)
	res := (rand8 < under)
	payout := bet * 98.5 / float64(under)
	r := &Roll{
		Key:        k.Ak,
		KeyObj:     k,
		Bet:        bet,
		RollUnder:  under,
		RollRandom: rand8,
		Payout:     payout,
		Result:     res,
		Before:     k.Balance,
		After:      0,
	}
	if res {
		r.After = k.Balance - bet + payout
	} else {
		r.After = k.Balance - bet
	}
	k.Balance = r.After
	k.TxCount++
	if k.Balance > k.Peak {
		k.Peak = k.Balance
	}
	if k.Balance < 1 {
		k.BankruptAt = &JsonTime{time.Now()}
	}
	tx := Gdb.Begin()
	dbResult := tx.Create(r)
	if dbResult.Error != nil {
		tx.Rollback()
		return nil
	}
	dbResult = tx.Save(k)
	if dbResult.Error != nil {
		tx.Rollback()
		return nil
	}
	tx.Commit()
	return r
}
