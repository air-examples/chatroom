package gas

import "github.com/sheng/air"

func PreRenderHandler() air.Gas {
	return func(next air.Handler) air.Handler {
		return func(req *air.Request, res *air.Response) error {
			if _, ok := req.Values["title"]; !ok {
				req.Values["title"] =
					req.LocalizedString("title")
			}
			return next(req, res)
		}
	}
}
