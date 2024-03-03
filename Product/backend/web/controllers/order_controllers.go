package controllers

import (
	"strconv"

	"github.com/Ferriem/go-web/Product/common"
	"github.com/Ferriem/go-web/Product/datamodels"
	"github.com/Ferriem/go-web/Product/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type OrderController struct {
	Ctx          iris.Context
	OrderService services.IOrderService
}

func (o *OrderController) Get() mvc.View {
	orderArray, err := o.OrderService.GetAllOrderInfo()
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	return mvc.View{
		Name: "order/view.html",
		Data: iris.Map{
			"order": orderArray,
		},
	}
}

func (o *OrderController) GetDelete() {
	idString := o.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	isOk := o.OrderService.DeleteOrderById(id)
	if !isOk {
		o.Ctx.Application().Logger().Debug("Delete order failed!")
	} else {
		o.Ctx.Application().Logger().Debug("Delete order success!")
	}
	o.Ctx.Redirect("/order")
}

func (o *OrderController) PostUpdate() {
	order := &datamodels.Order{}
	o.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "ferriem"})
	if err := dec.Decode(o.Ctx.Request().Form, order); err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}

	err := o.OrderService.UpdateOrder(order)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}

	o.Ctx.Redirect("/order")
}
