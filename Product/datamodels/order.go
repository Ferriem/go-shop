package datamodels

type Order struct {
	ID          int64 `sql:"ID" ferriem:"ID"`
	UserId      int64 `sql:"userId" ferriem:"UserID"`
	ProductId   int64 `sql:"productId" ferriem:"ProductID"`
	OrderStatus int   `sql:"orderStatus" ferriem:"OrderStatus"`
	TrackNum    int64 `sql:"trackNum" ferriem:"TrackNum"`
}

const (
	OrderWait = iota
	OrderSuccess
	OrderFailed
)
