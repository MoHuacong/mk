package push

import (
	"github.com/MoHuacong/wic"
)

type JsonType struct {
	Type string `json:"type"`
}

type JsonUserInfo struct {
	Url string `json:"url"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type JsonData struct {
	List []uint64 `json:"list"`
	Data string `json:"data"`
}

type JsonClient struct {
	Url string `json:"url"`
}

type JsonRequest struct {
	JsonType
	Id uint64 `json:"id"`
	Data string `json:"data"`
}

type JsonResponse struct {
	JsonType
	JsonUserInfo
	Id JsonData `json:"id"`
}

type UserStorage struct {
	JsonUserInfo
	Url string `json:"url"`
}

type UserIdStorage struct {
	UserStorage
	Id map[*wic.Fd]bool
}

type Users map[string]*UserIdStorage