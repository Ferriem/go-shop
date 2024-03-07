package services

import (
	"github.com/Ferriem/go-web/Product/datamodels"
	"github.com/Ferriem/go-web/Product/repositories"
)

type IProductService interface {
	GetProductByID(int64) (*datamodels.Product, error)
	GetAllProduct() ([]*datamodels.Product, error)
	DeleteProductByID(int64) bool
	InsertProduct(product *datamodels.Product) (int64, error)
	UpdateProduct(product *datamodels.Product) error
	SubNumberOne(int64) error
}

type ProductService struct {
	productRepository repositories.IProduct
}

func NewProductService(repository repositories.IProduct) IProductService {
	return &ProductService{productRepository: repository}
}

func (p *ProductService) GetProductByID(productID int64) (*datamodels.Product, error) {
	return p.productRepository.SelectByKey(productID)
}

func (p *ProductService) GetAllProduct() ([]*datamodels.Product, error) {
	return p.productRepository.SelectAll()
}

func (p *ProductService) DeleteProductByID(productID int64) bool {
	return p.productRepository.Delete(productID)
}

func (p *ProductService) InsertProduct(product *datamodels.Product) (productID int64, err error) {
	return p.productRepository.Insert(product)
}

func (p *ProductService) UpdateProduct(product *datamodels.Product) (err error) {
	return p.productRepository.Update(product)
}

func (p *ProductService) SubNumberOne(productID int64) error {
	return p.productRepository.SubProductByOne(productID)
}
