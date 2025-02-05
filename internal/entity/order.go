package entity

type OrderCreateReq struct {
	UserID    string `json:"user_id"`
}

type OrderGetReq struct {
	ID string `json:"id"`
}

type OrderGetRes struct {
	Order Order `json:"order"`
}


type OrderUpt struct {
	ID   string  `json:"id"`
	Type      string `json:"type"`
	Quantity  int32  `json:"quantity"`
	TotalPrice float32 `json:"total_price"`
	Status    string `json:"status"`
}

type OrderDeleteReq struct {
	ID string `json:"id"`
}

type OrderListsReq struct {
	UserID  string     `json:"user_id"`
	Filter  Pagination `json:"filter"`
	Prtype    string   `json:"type"`
}

type OrderListsRes struct {
	Orders      []Order   `json:"orders"`
	Pagination  Pagination `json:"pagination"`
	TotalCount  int      `json:"total_count"`
}

type Order struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	ItemID     string `json:"item_id"`
	Type       string `json:"type"`
	Quantity   int32  `json:"quantity"`
	TotalPrice string `json:"total_price"`
	Status     string `json:"status"`
	CreatedAt  string
	UpdatedAt  string
}



