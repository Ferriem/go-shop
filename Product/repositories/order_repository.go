package repositories

import (
	"database/sql"
	"strconv"

	"github.com/Ferriem/go-web/Product/common"
	"github.com/Ferriem/go-web/Product/datamodels"
)

var Database = "product"

type IOrder interface {
	Conn() error
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

type OrderManager struct {
	table     string
	mysqlConn *sql.DB
}

func NewOrderManager(table string, DB *sql.DB) IOrder {
	return &OrderManager{table: table, mysqlConn: DB}
}

func (o *OrderManager) Conn() error {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = "`order`"
	}
	return nil
}

func (o *OrderManager) Insert(order *datamodels.Order) (productid int64, err error) {
	if err = o.Conn(); err != nil {
		return
	}
	sql := "INSERT " + o.table + " SET userId=?,productId=?,orderStatus=?,trackNum=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return 0, err
	}
	result, err := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus, order.TrackNum)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (o *OrderManager) Delete(productId int64) bool {
	if err := o.Conn(); err != nil {
		return false
	}

	sql := "DELETE FROM " + o.table + " WHERE ID = ?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}
	_, err = stmt.Exec(productId)
	if err != nil {
		return false
	}
	return true
}

func (o *OrderManager) Update(order *datamodels.Order) error {
	if err := o.Conn(); err != nil {
		return err
	}

	sql := "UPDATE " + o.table + " SET userId=?,productId=?,orderStatus=?,trackNum=? WHERE ID = ?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.UserId, order.ProductId, order.OrderStatus, order.TrackNum, order.ID)
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderManager) SelectByKey(productId int64) (*datamodels.Order, error) {
	if err := o.Conn(); err != nil {
		return &datamodels.Order{}, err
	}
	sql := "SELECT * FROM " + o.table + " WHERE ID = " + strconv.FormatInt(productId, 10)
	row, err := o.mysqlConn.Query(sql)
	defer row.Close()
	if err != nil {
		return &datamodels.Order{}, err
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, nil
	}
	order := &datamodels.Order{}
	common.DataToStructByTagSql(result, order)
	return order, nil
}

func (o *OrderManager) SelectAll() (orderArray []*datamodels.Order, err error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}
	sql := "SELECT * FROM " + o.table
	rows, err := o.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	result := common.GetResultRows(rows)
	if len(result) == 0 {

		return nil, nil
	}

	for _, v := range result {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		orderArray = append(orderArray, order)
	}
	return
}

func (o *OrderManager) SelectAllWithInfo() (map[int]map[string]string, error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}
	sql := "SELECT o.ID,p.productName,o.orderStatus,o.trackNum FROM " + Database + "." + o.table + " as o left join product as p on o.productId = p.ID"
	rows, err := o.mysqlConn.Query(sql)
	if err != nil {
		return nil, err
	}
	return common.GetResultRows(rows), nil
}
