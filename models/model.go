package models

import (
	"encoding/base64"

	cmap "github.com/orcaman/concurrent-map"
)

var (
	NameMap cmap.ConcurrentMap
)

type UserInfo struct {
	Name string
}

func InitModel() {
	NameMap = cmap.New()
}

func GetAuthKey(name string) string {
	return base64.StdEncoding.EncodeToString([]byte(name))
}

func NewUser(name string) UserInfo {
	return UserInfo{
		Name: name,
	}
}

func GetUser(key string) (UserInfo, bool) {
	v, ok := NameMap.Get(key)
	if !ok {
		return UserInfo{}, false
	}
	res, ok := v.(UserInfo)
	return res, ok
}

func SetUser(u UserInfo) {
	NameMap.Set(GetAuthKey(u.Name), u)
}
