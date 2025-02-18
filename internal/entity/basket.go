package entity


type BasketItem struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	UserId    string  `json:"user_id"`
	Price     float64 `json:"price"`
	Count     int     `json:"count"`
	Status    string  `json:"status"` // e.g., "sold", "not_sold"
}


type BasketResponse struct{
	Id  	    []string 	
	TotalPrice  float64 `json:"price"`
	Count       int     `json:"count"`
}


type BasketItemForSwagger struct {
	ProductID string  `json:"product_id"`
	Count     int	  `json:"count"`
}

type ListItem struct {
	ID       string   `json:"id"`
	Price    float64  `json:"price"`
	Count    int      `json:"count"`
	Product  Product  `json:"product"`
	Pictures []string `json:"pictures"`
}


type ListBasketItem struct {
	Items      []ListItem `json:"items"`
	TotalPrice float64    `json:"total_price"`
	TotalCount int        `json:"total_count"`
}

type BasketDelete struct{
	Basketid string
	Userid string
}


type BasketProductPrice struct{
	Price          float64  `json:"price"`
	DiscountPrice  float64  `json:"discount_price"`
	FinalPrice     float64  `json:"final_price"`
}