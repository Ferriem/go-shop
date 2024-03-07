package main

import (
	"fmt"

	"github.com/Ferriem/go-web/Product/common"
	"github.com/Ferriem/go-web/Product/repositories"
	"github.com/Ferriem/go-web/Product/services"
	"github.com/Ferriem/go-web/RabbitMQ"
)

func main() {
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}

	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderManager("`order`", db)
	orderService := services.NewOrderService(order)

	rabbitmqConsumeSimple := RabbitMQ.NewRabbitMQSimple("product")
	rabbitmqConsumeSimple.ConsumeSimple(orderService, productService)
}
