package model

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/atomic"

	"time"
)

var Stat = &struct {
	Uptime       time.Time
	RequestCount atomic.Int64
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
		Stat.RequestCount.Inc()
		return nil
	}
}
