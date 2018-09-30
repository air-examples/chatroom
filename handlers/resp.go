package handlers

import "github.com/aofei/air"

func Error(code int, err error) error {
	return &air.Error{
		Code:    code,
		Message: err.Error(),
	}
}

func Success(res *air.Response, data interface{}) error {
	ret := map[string]interface{}{}
	ret["code"] = 0
	ret["error"] = ""
	ret["data"] = data
	if data == nil {
		ret["data"] = ""
	}
	return res.JSON(ret)
}
