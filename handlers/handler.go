package handlers

import (
	"github.com/air-examples/chatroom/gas"
	"github.com/air-examples/chatroom/utils"
	"github.com/sheng/air"
)

func init() {
	air.ErrorHandler = errorHandler
	air.GET("/", indexHandler, gas.PreRender)
}

func errorHandler(err error, req *air.Request, res *air.Response) {
	e, ok := err.(*air.Error)
	if !ok {
		e.Code = 500
		e.Message = "Server Internal Error"
	}
	if !res.Written {
		if req.Method == "GET" || req.Method == "HEAD" {
			delete(res.Headers, "ETag")
			delete(res.Headers, "Last-Modified")
		}
		ret := utils.M{}
		ret["data"] = ""
		ret["code"] = e.Code
		ret["error"] = e.Message
		res.StatusCode = e.Code
		if e.Code != 500 {
			res.JSON(ret)
			return
		}
		for k, v := range ret {
			req.Values[k] = v
		}
		res.Render(req.Values, "error.html")
		return
	}

}

func indexHandler(req *air.Request, res *air.Response) error {
	return res.Render(req.Values, "index.html", "base.html")
}
