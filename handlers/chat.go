package handlers

import (
	"github.com/air-examples/chatroom/models"

	"github.com/aofei/air"
	cmap "github.com/orcaman/concurrent-map"
)

var (
	users = cmap.New()
)

func init() {
	air.GET("/chat", chatPageHandler, Auth)
	air.POST("/api/name", chatNameHandler)
	air.GET("/socket", socketHandler, Auth)
}

func chatPageHandler(req *air.Request, res *air.Response) error {
	return res.Render(req.Values, "chat.html", "base.html")
}

func chatNameHandler(req *air.Request, res *air.Response) error {
	nameT := req.Params["name"].FirstValue()
	if nameT == nil {
		air.ERROR("param name missed")
		return Error(
			res,
			400,
			"bad request (miss param `name`)",
		)
	}
	name := nameT.String()
	k := models.GetAuthKey(name)
	if _, ok := models.GetUser(k); ok {
		air.ERROR("duplicate name", Map{
			"name":    name,
			"request": req,
		})
		return Error(res, 400, "duplicate name")
	}

	models.SetUser(models.NewUser(name))
	res.Cookies["name"] = &air.Cookie{
		Name:  "name",
		Value: k,
		Path:  "/",
	}

	return Success(res, "")
}

func socketHandler(req *air.Request, res *air.Response) error {
	c, err := res.WebSocket()
	if err != nil {
		air.ERROR("upgrade to websocket error", Map{
			"request": req,
			"error":   err.Error(),
		})
		return Error(res, 500, err.Error())
	}

	defer c.Close()

	name, _ := req.Values["name"].(string)
	me := newSocketManager(name)
	if _, ok := users.Get(name); !ok {
		users.Set(name, me)
	}

	c.TextHandler = func(text string) error {
		me.SendMsg(newMsg(name, text))
		return nil
	}

	c.ErrorHandler = func(err error) {
		air.ERROR("websocket error", Map{
			"err": err.Error(),
		})
		me.Close()
	}

	for {
		select {
		case <-me.newMsg:
			if text, err := me.msg.Marshal(); err == nil {
				err = c.WriteText(string(text))
				if err != nil {
					air.ERROR("send socket msg error",
						Map{
							"content": me.msg,
							"to":      me.name,
							"err":     err,
						})
					me.Close()
				}
				<-me.writeChan
			}

		case <-me.shutdown:
			break
		}

	}

	return nil
}
