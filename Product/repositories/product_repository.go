package repositories

import (
	"database/sql"
	"strconv"

	"github.com/Ferriem/go-web/Product/common"
	"github.com/Ferriem/go-web/Product/datamodels"
)

//define the interface

type IProduct interface {
	Conn() error
	Insert(product *datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(product *datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
}

//implement the interface

type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

func NewProductManager(table string, DB *sql.DB) IProduct {
	return &ProductManager{table: table, mysqlConn: DB}
}

// connect to the database
func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}
	if p.table == "" {
		p.table = "product"
	}
	return
}

// insert a product
func (p *ProductManager) Insert(product *datamodels.Product) (productid int64, err error) {
	if err = p.Conn(); err != nil {
		return
	}

	sql := "INSERT " + p.table + " SET productName=?,productNum=?,productImage=?,productUrl=?,productPrice=?,productDiscount=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return 0, err
	}
	result, err := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl, product.ProductPrice, product.ProductDiscount)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()

}

// delete a product
func (p *ProductManager) Delete(productID int64) bool {
	if err := p.Conn(); err != nil {
		return false
	}

	sql := "DELETE FROM " + p.table + " WHERE ID = ?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}
	_, err = stmt.Exec(productID)
	if err != nil {
		return false
	}
	return true
}

// update a product
func (p *ProductManager) Update(product *datamodels.Product) error {
	if err := p.Conn(); err != nil {
		return err
	}

	sql := "UPDATE " + p.table + " SET productName=?,productNum=?,productImage=?,productUrl=?,productPrice=?,productDiscount=? WHERE ID = ?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl, product.ProductPrice, product.ProductDiscount, product.ID)
	if err != nil {
		return err
	}
	return nil
}

// select a product by primary key
func (p *ProductManager) SelectByKey(productID int64) (*datamodels.Product, error) {
	if err := p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}

	sql := "SELECT * FROM " + p.table + " WHERE ID = " + strconv.FormatInt(productID, 10)
	row, err := p.mysqlConn.Query(sql)
	defer row.Close()
	if err != nil {
		return &datamodels.Product{}, err
	}

	result := common.GetResultRow(row)

	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}

	productResult := &datamodels.Product{}
	common.DataToStructByTagSql(result, productResult)
	return productResult, nil

}

// select all products
func (p *ProductManager) SelectAll() (productArray []*datamodels.Product, err error) {
	if err := p.Conn(); err != nil {
		return nil, err
	}

	sql := "SELECT * FROM " + p.table
	rows, err := p.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, nil
	}

	for _, v := range result {
		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)
		productArray = append(productArray, product)
	}
	return
}
