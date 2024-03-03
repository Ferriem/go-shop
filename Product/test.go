package main

import (
	"github.com/Ferriem/go-web/Product/common"
	"github.com/Ferriem/go-web/Product/datamodels"
)

func main() {
	data := map[string]string{
		"ID":           "1",
		"productName":  "test",
		"productNum":   "100",
		"productImage": "test.jpg",
		"productUrl":   "www.test.com",
	}
	product := &datamodels.Product{}
	common.DataToStructByTagSql(data, product)
}
