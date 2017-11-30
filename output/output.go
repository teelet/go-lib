package output

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"sdk.look.360.cn/application/library/app"
	"sdk.look.360.cn/application/library/badcase"
)

type JS map[string]interface{}

func Auto(data interface{}, c echo.Context) error {
	return AutoExtra(data, JS{}, c)
}

func AutoExtra(data interface{}, extra map[string]interface{}, c echo.Context) error {
	isMap := true
	mapData := JS{}
	arrayData := []interface{}{}

	errMessage := ""
	errCode := 0
	errBacktrace := ""

	switch data.(type) {
	case error:
		if tempCase, ok := data.(*badcase.Stack); ok {
			errCode = tempCase.Code()
			errMessage = tempCase.Error()
			errBacktrace = tempCase.Backtrace()
		} else {
			tempError := data.(error)
			errCode = 1
			errMessage = tempError.Error()
		}
	case map[string]interface{}:
		mapData = data.(map[string]interface{})
		isMap = true
	case JS:
		mapData = data.(JS)
		isMap = true
	case []interface{}:
		arrayData = data.([]interface{})
		isMap = false
	}

	f := c.QueryParam("f")
	cb := c.QueryParam("callback") // JSONP Vulnerability [http://blog.knownsec.com/2015/03/jsonp_security_technic/]

	if cb == "" {
		cb = "undefined_callback"
	}

	response := JS{}
	if _, ok := mapData["errno"]; ok {
		response = mapData
	} else {
		response = JS{
			"errno":  errCode,
			"errmsg": errMessage,
			"data":   map[string]interface{}{},
		}
		if isMap && len(mapData) > 0 {
			response["data"] = mapData
		} else if len(arrayData) > 0 {
			response["data"] = arrayData
		}
	}

	for k, v := range extra {
		response[k] = v
	}

	if f == "jsonp" {
		c.JSONP(http.StatusOK, cb, response)
	} else {
		c.JSON(http.StatusOK, response)
	}

	if errCode > 0 {
		if errBacktrace != "" {
			app.LogAlarm(c, errBacktrace, errCode)
		} else {
			app.LogAlarm(c, fmt.Sprintf("code=%d %s", errCode, errMessage), errCode)
		}
	}

	return nil
}

func Error(msg string, code int, c echo.Context) error {
	return AutoExtra(JS{
		"errno":  code,
		"errmsg": msg,
	}, JS{}, c)
}

func ErrorExtra(msg string, code int, extra map[string]interface{}, c echo.Context) error {
	return AutoExtra(JS{
		"errno":  code,
		"errmsg": msg,
	}, extra, c)
}
