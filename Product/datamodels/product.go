package datamodels

type Product struct {
	ID           int64  `json:"id" sql:"ID" ferriem:"ID"`
	ProductName  string `json:"ProductName" sql:"productName" ferriem:"productName"`
	ProductNum   int64  `json:"ProductNum" sql:"productNum" ferriem:"productNum"`
	ProductImage string `json:"ProductImage" sql:"productImage" ferriem:"productImage"`
	ProductUrl   string `json:"ProductUrl" sql:"productUrl" ferriem:"productUrl"`
}
