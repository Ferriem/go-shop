package middleware

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/Ferriem/go-web/Product/common"
	"github.com/Ferriem/go-web/Product/datamodels"
	"github.com/Ferriem/go-web/Product/encrypt"
	"github.com/kataras/iris/v12"
)

func AuthConnProduct(ctx iris.Context) {
	sign := ctx.GetCookie("sign")
	if sign == "" {
		ctx.Application().Logger().Debug("sign is empty")
		ctx.Redirect("/user/login")
		return
	}
	result, err := encrypt.DePwdCode(sign)
	if err != nil {
		ctx.Application().Logger().Debug(err)
	}
	userInfo := datamodels.UserInfo{}
	err = json.Unmarshal(result, &userInfo)
	if err != nil {
		ctx.Application().Logger().Debug(err)
	}
	rdb := common.NewRedisConn()
	ctxs := context.Background()
	results, err := rdb.HVals(ctxs, userInfo.UserName).Result()
	if err != nil {
		ctx.Application().Logger().Debug(err)
	}
	if len(results) == 0 {
		ctx.Application().Logger().Debug("results is empty")
		ctx.Redirect("/user/login")
		return
	}
	userID, err := strconv.Atoi(results[0])
	if userInfo.ID != int64(userID) {
		ctx.Application().Logger().Debug("userInfo.ID != userID")
		ctx.Redirect("/user/login")
		return
	}
	ctx.Application().Logger().Debug("Already login")
	ctx.Next()
}
