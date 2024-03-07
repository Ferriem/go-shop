package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Ferriem/go-web/Product/common"
	"github.com/Ferriem/go-web/Product/datamodels"
	"github.com/Ferriem/go-web/Product/encrypt"
	"github.com/Ferriem/go-web/Product/services"
	"github.com/Ferriem/go-web/Product/tool"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type UserController struct {
	Ctx     iris.Context
	Service services.IUserService
}

func (c *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "user/register.html",
	}
}

func (c *UserController) PostRegister() {
	// var (
	// 	nickName = c.Ctx.FormValue("nickName")
	// 	userName = c.Ctx.FormValue("userName")
	// 	password = c.Ctx.FormValue("password")
	// )

	user := &datamodels.User{}
	c.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "form"})
	if err := dec.Decode(c.Ctx.Request().Form, user); err != nil {
		c.Ctx.Application().Logger().Debug(err)
	}
	_, err := c.Service.AddUser(user)
	if err != nil {
		c.Ctx.Redirect("/user/error")
	}
	c.Ctx.Redirect("/user/login")
	return
}

func (c *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name: "user/login.html",
	}
}

func (c *UserController) PostLogin() mvc.Response {
	user := &datamodels.User{}
	c.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "form"})
	if err := dec.Decode(c.Ctx.Request().Form, user); err != nil {
		c.Ctx.Application().Logger().Debug(err)
	}
	user, isOk := c.Service.IsPwdSuccess(user.UserName, user.HashPassword)
	if !isOk {
		return mvc.Response{
			Path: "/user/login",
		}
	}
	userInfo := &datamodels.UserInfo{
		ID:       user.ID,
		UserName: user.UserName,
		IP:       c.Ctx.Request().RemoteAddr,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		c.Ctx.Application().Logger().Debug(err)
	}

	uidString, err := encrypt.EnPwdCode(jsonData)
	if err != nil {
		fmt.Println(err)
	}
	tool.GlobalCookie(c.Ctx, "sign", uidString)
	tool.GlobalCookie(c.Ctx, "uid", strconv.FormatInt(user.ID, 10))

	return mvc.Response{
		Path: "/product/",
	}
}
