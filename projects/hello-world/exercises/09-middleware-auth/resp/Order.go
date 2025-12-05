package resp

type Order struct {
	Id       uint64 `json:"id"`
	Products []Product
}

type Product struct {
	Id          uint64 `json:"id"`
	ProductName string `json:"productName"`
}
