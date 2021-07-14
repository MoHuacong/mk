package push

import (
	"fmt"
	"github.com/MoHuacong/wic/tools"
)

type Api struct {
	MySQL
	err *tools.Stack
}

func NewApi() *Api {
	return &Api{err: tools.NewStack()}
}

/* 注册 */
func (api *Api) Register(user, pass, url string) uint {
	var userData User
	rows, err := api.db.Query("SELECT * from user WHERE user='"+user+"' or url='"+url+"'")
	if err != nil {
		return api.formatExit(1,"sql命令错误1")
	}
	
	for rows.Next() {
		api.Scan(rows, &userData)
	}
	
	if userData.User != "" {
		return api.formatExit(2, "帐号或密码URL已存在")
	}
	
	sql := fmt.Sprintf("INSERT INTO user (id, user, pass, url, date) VALUES (NULL, '%s', '%s', '%s', now())", user, pass, url)
	if _, err = api.db.Query(sql); err != nil {
		return api.formatExit(3, "注册失败")
	}
	
	return api.formatExit(0, "注册成功")
}

/* 登录 */
func (api *Api) Login(user, pass string) uint {
	var userData User
	rows, err := api.db.Query("SELECT * from user WHERE user='"+user+"' and pass='"+pass+"'")
	if err != nil {
		return api.formatExit(1,"sql命令错误1")
	}

	for rows.Next() {
		api.Scan(rows, &userData)
	}

	if userData.User == "" {
		return api.formatExit(2, "帐号或密码不正确")
	}
	return api.formatExit(0, "成功")
}

/* 判断url存在 */
func (api *Api) IsUrl(user *User, url string) uint {
	rows, _ := api.db.Query("SELECT * from user WHERE url='"+url+"'")
	for rows.Next() {
		api.Scan(rows, user)
	}
	if user.Url == "" {
		return api.formatExit(1, "url不存在")
	}
	return api.formatExit(0, "OK")
}

func (api *Api) Is(user, pass, url string) uint {
	var userData User
	rows, err := api.db.Query("SELECT * from user WHERE user='"+user+"' and pass='"+pass+"' and url='"+url+"'")
	if err != nil {
		return api.formatExit(1,"sql命令错误1")
	}

	for rows.Next() {
		api.Scan(rows, &userData)
	}
	if userData.Url == "" {
		return api.formatExit(2, "url不存在")
	}
	return api.formatExit(0, "OK")
}

func (api *Api) formatExit(ret uint, v interface{}) uint {
	api.err.Push(v)
	return ret
}

func (api *Api) Error() string {
	data := api.err.Pop()
	if data == nil {
		return ""
	}
	return data.(string)
}

func (api *Api) String() string {
	return api.Error()
}