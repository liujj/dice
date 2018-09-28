package handler

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"net/http"
	"strconv"
	"tigerMachine/model"
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
	return c.Redirect(302,"/static/index.html")
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
	betFloat, err := strconv.ParseFloat(bet, 32)

	r := model.MakeRoll(k, float32(betFloat), uint8(underInt))
	DataResponse(c, http.StatusOK, 0, "ok", struct {
		K *model.Key  `json:"k"`
		R *model.Roll `json:"r"`
	}{k, r})
	return nil
}

func Status(c echo.Context) (err error) {
	k, err := getKey(c)

	if err != nil {
		return SimpleResponse(c, 404, -1, err.Error())
	} else if k == nil {
		return SimpleResponse(c, 404, -1, "internal error")
	}

	DataResponse(c, http.StatusOK, 0, "ok", struct {
		K *model.Key `json:"k"`
	}{k})
	return nil
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
