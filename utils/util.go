package utils

import (
	"encoding/base64"

	"github.com/aofei/air"
)

type M map[string]interface{}

func Error(code int, err error) error {
	return &air.Error{
		Code:    code,
		Message: err.Error(),
	}
}

func Success(res *air.Response, data interface{}) error {
	ret := M{}
	ret["code"] = 0
	ret["error"] = ""
	ret["data"] = data
	if data == nil {
		ret["data"] = ""
	}
	return res.JSON(ret)
}

func Base64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
