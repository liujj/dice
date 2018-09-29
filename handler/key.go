package handler

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"tigerMachine/model"
)

func NewKey(c echo.Context) error {
	u, err := model.NewKey()
	if err != nil {
		return SimpleResponse(c, http.StatusBadRequest, -1, err.Error())
	} else {
		return DataResponse(c, http.StatusOK, 0, "ok", struct {
			*model.Key
			Url string `json:"url""`
		}{u, c.Request().Host + "/start?ak=" + u.Ak + "&sk=" + u.Sk})
	}
}

func QueryKey(c echo.Context) error {
	ak := c.Param("key")
	if ak == "" {
		return SimpleResponse(c, http.StatusNotFound, -1, "no Access Key input")
	}
	u, err := model.QueryKey(ak)
	if err != nil {
		return SimpleResponse(c, http.StatusBadRequest, -1, err.Error())
	} else {
		return DataResponse(c, http.StatusOK, 0, "ok", u)
	}
}

func ListKeys(c echo.Context) error {
	frInt, _ := strconv.Atoi(c.QueryParam("fr"))
	toInt, _ := strconv.Atoi(c.QueryParam("to"))
	showDetail := c.QueryParam("detail")
	if toInt == 0 {
		toInt = frInt + 50
	}
	if frInt == 0 {
		frInt = 1
	}
	frInt = frInt - 1
	if frInt < 0 || toInt < 0 || frInt > toInt {
		return SimpleResponse(c, http.StatusBadRequest, -1, "page info error")
	}
	keys, err := model.ListKeys(frInt, toInt)
	if err != nil {
		return SimpleResponse(c, http.StatusBadRequest, -1, err.Error())
	}
	if len(keys) == 0 {
		return SimpleResponse(c, http.StatusNotFound, -1, "zero item found")
	}
	if showDetail == "true" {
		return c.JSON(http.StatusOK, &struct {
			ResponseBase
			ResponsePage `json:"page"`
			Data         interface{} `json:"data"`
		}{ResponseBase: ResponseBase{Code: 0, Msg: "ok"}, ResponsePage: ResponsePage{frInt + 1, frInt + len(keys)}, Data: keys})
	} else {
		keySummary := []string{}
		for _, k := range keys {
			keySummary = append(keySummary, k.Ak)
		}
		return c.JSON(http.StatusOK, &struct {
			ResponseBase
			ResponsePage `json:"page"`
			Data         interface{} `json:"data"`
		}{ResponseBase: ResponseBase{Code: 0, Msg: "ok"}, ResponsePage: ResponsePage{frInt + 1, frInt + len(keys)}, Data: keySummary})
	}

}

func SetKey(c echo.Context) error {
	ak := c.Param("key")
	params := c.QueryParams()
	param := make(map[string]string)
	for _, colEnum := range []string{"currency", "ip", "vaild", "freq", "remark", "expire", "expireAt"} {
		if value, exist := params[colEnum]; exist {
			param[colEnum] = value[0]
		}
	}
	if len(param) == 0 || ak == "" {
		return SimpleResponse(c, http.StatusNotFound, -1, "Input error")
	}
	err := model.SetKey(ak, param)
	if err != nil {
		return SimpleResponse(c, http.StatusBadRequest, -1, err.Error())
	} else {
		return SimpleResponse(c, http.StatusOK, 0, "Modified")
	}
}

func DelKey(c echo.Context) error {
	ak := c.Param("key")
	if ak == "" {
		return SimpleResponse(c, http.StatusNotFound, -1, "no Access Key input")
	}
	err := model.DelKey(ak)
	if err != nil {
		return SimpleResponse(c, http.StatusBadRequest, -1, err.Error())
	} else {
		return SimpleResponse(c, http.StatusOK, 0, "Deleted")
	}
}

func QueryKeybyCache(c echo.Context) error {
	ak := c.Param("key")
	model.KC.Mutex.RLock()
	defer model.KC.Mutex.RUnlock()
	if c.Path() == "/key/cache/list" {
		mapList := []string{}
		for keyAk := range model.KC.Map {
			mapList = append(mapList, keyAk)
		}
		return DataResponse(c, http.StatusOK, 0, "ok", mapList)
	} else if ak == "" {
		return SimpleResponse(c, http.StatusNotFound, -1, "no Access Key input")
	} else {
		return DataResponse(c, http.StatusOK, 0, "ok", model.KC.Map[ak])
	}
}

func ResetKeyCache(c echo.Context) error {
	model.ClearKeyCache()
	return SimpleResponse(c, http.StatusOK, 0, "ok")
}
