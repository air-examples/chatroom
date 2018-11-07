package handlers

import (
	"errors"

	"github.com/aofei/air"
)

func init() {
	air.ErrorHandler = errorHandler
	air.STATIC("/assets", air.AssetRoot)
	air.GET("/", indexHandler)
	air.GET("/api/const", constHandler)
}

type Map map[string]interface{}

func errorHandler(err error, req *air.Request, res *air.Response) {
	if res.Written {
		return
	}

	air.ERROR(err.Error())

	if req.Method == "GET" || req.Method == "HEAD" {
		delete(res.Headers, "ETag")
		delete(res.Headers, "Last-Modified")
		if req.Method == "GET" && res.Status == 401 {
			res.Status = 307
			res.Redirect("http://" + air.Address)
			return
		}
	}

	res.WriteJSON(Map{
		"data":  "",
		"code":  res.Status,
		"error": err.Error(),
	})

}

func indexHandler(req *air.Request, res *air.Response) error {
	return res.Render(req.Values, "index.html", "base.html")
}

func constHandler(req *air.Request, res *air.Response) error {
	c := Map{
		"message can not be empty":               "",
		"join chatroom":                          "",
		"name repeat, please input anothor name": "",
		"connection closed":                      "",
	}
	for k, _ := range c {
		c[k] = req.LocalizedString(k)
	}
	return Success(res, c)
}

func Error(res *air.Response, code int, message string) error {
	res.Status = code
	return errors.New(message)
}

func Success(res *air.Response, data interface{}) error {
	res.Status = 200
	if data == nil {
		data = ""
	}
	return res.WriteJSON(Map{
		"code":  0,
		"error": "",
		"data":  data,
	})
}
