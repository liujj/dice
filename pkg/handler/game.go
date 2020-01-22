package handler

import (
	"dice/pkg/model"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Start(c echo.Context) (err error) {
	ak := c.QueryParam("ak")
	sk := c.QueryParam("sk")
	if len(ak) != 8 {
		return SimpleResponse(c, http.StatusForbidden, -1, "ak/sk not match")
	}
	if len(sk) != 32 {
		return SimpleResponse(c, http.StatusForbidden, -1, "ak/sk not match")
	}
	k, err := model.QueryKey(ak)
	if err != nil {
		return SimpleResponse(c, http.StatusForbidden, -1, "ak/sk not match")
	}
	if k.Sk != sk {
		return SimpleResponse(c, http.StatusForbidden, -1, "ak/sk not match")
	}
	sess, _ := session.Get("pubilc_session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
	}
	sess.Values["ak"] = ak
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusFound, "/static/index.html")
}

func Roll(c echo.Context) (err error) {
	k, err := getKey(c)

	if err != nil {
		return SimpleResponse(c, 404, -1, err.Error())
	} else if k == nil {
		return SimpleResponse(c, 404, -1, "internal error")
	}

	under := c.QueryParam("under")
	underInt, err := strconv.Atoi(under)
	bet := c.QueryParam("bet")
	betFloat, err := strconv.ParseFloat(bet, 64)
	betFloat = math.Floor(betFloat*100) / 100

	r := model.MakeRoll(k, float64(betFloat), uint8(underInt))
	return DataResponse(c, http.StatusOK, 0, "ok", struct {
		K *model.Key  `json:"k"`
		R *model.Roll `json:"r"`
	}{k, r})
}

func Status(c echo.Context) (err error) {
	k, err := getKey(c)

	if err != nil {
		return SimpleResponse(c, 404, -1, err.Error())
	} else if k == nil {
		return SimpleResponse(c, 404, -1, "internal error")
	}

	rs := k.GetHistory(0, 25)

	return DataResponse(c, http.StatusOK, 0, "ok", struct {
		K  *model.Key   `json:"k"`
		Rs []model.Roll `json:"rs"`
	}{k, rs})
}

func LastRecord(c echo.Context) (err error) {
	return nil
}

func getKey(c echo.Context) (*model.Key, error) {
	sess, _ := session.Get("pubilc_session", c)
	ak := sess.Values["ak"].(string)
	k, err := model.FetchKey(ak)
	return k, err
}
