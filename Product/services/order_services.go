package services

import (
	"github.com/Ferriem/go-web/Product/datamodels"
	"github.com/Ferriem/go-web/Product/repositories"
)

type IOrderService interface {
	GetOrderById(int64) (*datamodels.Order, error)
	GetAllOrder() ([]*datamodels.Order, error)
	DeleteOrderById(int64) bool
	InsertOrder(order *datamodels.Order) (int64, error)
	UpdateOrder(order *datamodels.Order) error
	GetAllOrderInfo() (map[int]map[string]string, error)
	InsertOrderByMessage(message *datamodels.Message) (int64, error)
}

type OrderService struct {
	orderRepository repositories.IOrder
}

func NewOrderService(repository repositories.IOrder) IOrderService {
	return &OrderService{orderRepository: repository}
}

func (o *OrderService) GetOrderById(orderId int64) (*datamodels.Order, error) {
	return o.orderRepository.SelectByKey(orderId)
}

func (o *OrderService) GetAllOrder() ([]*datamodels.Order, error) {
	return o.orderRepository.SelectAll()
}

func (o *OrderService) DeleteOrderById(orderId int64) bool {
	return o.orderRepository.Delete(orderId)
}

func (o *OrderService) InsertOrder(order *datamodels.Order) (int64, error) {
	return o.orderRepository.Insert(order)
}

func (o *OrderService) UpdateOrder(order *datamodels.Order) error {
	return o.orderRepository.Update(order)
}

func (o *OrderService) GetAllOrderInfo() (map[int]map[string]string, error) {
	return o.orderRepository.SelectAllWithInfo()
}

func (o *OrderService) InsertOrderByMessage(message *datamodels.Message) (int64, error) {
	order := &datamodels.Order{
		UserId:      message.UserID,
		ProductId:   message.ProductID,
		OrderStatus: datamodels.OrderSuccess,
	}
	return o.InsertOrder(order)
}
