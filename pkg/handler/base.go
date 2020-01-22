package handler

import (
	"github.com/labstack/echo/v4"
)

type ResponseBase struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type ResponsePage struct {
	FrItem int `json:"frItem"`
	ToItem int `json:"toItem"`
}

type Response struct {
	ResponseBase
	Data interface{} `json:"data"`
}

func SimpleResponse(c echo.Context, httpCode, code int, msg string) error {
	res := &ResponseBase{Code: code, Msg: msg}
	return c.JSON(httpCode, res)
}

func DataResponse(c echo.Context, httpCode, code int, msg string, data interface{}) error {
	res := &Response{ResponseBase: ResponseBase{Code: code, Msg: msg}, Data: data}
	return c.JSON(httpCode, res)
}
