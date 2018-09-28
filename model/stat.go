package model

import (
	"github.com/labstack/echo"
	"sync"
	"time"
)

var Stat = &struct {
	Uptime       time.Time
	RequestCount uint64
	Mutex        sync.RWMutex
}{
	Uptime: time.Now(),
}

func InitStat() {
	Stat.Uptime = time.Now().Local()
}

func StatMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		Stat.Mutex.Lock()
		defer Stat.Mutex.Unlock()
		Stat.RequestCount++
		return nil
	}
}
