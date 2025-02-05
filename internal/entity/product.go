package entity

type ProductCreate struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       float32 `json:"price"`
	SalePrice   float32 
	Color       string
	Size 		string
	PictureUrl  string `json:"picture_url"`
	Category_id string `json:"category_id"`
	PrType      string
}

type ProductCreateForSwagger struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       float32 `json:"price"`
	PictureUrl  string `json:"picture_url"`
	Category_id string `json:"category_id"`
	PrType      string
}

type ProductUpt struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}

type ProductPicture struct {
	ProductId  string `json:"product_id"`
	PictureUrl string `json:"picture_url"`
}

type ProductGet struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	PictureUrls []string `json:"picture_urls"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type ProductFilter struct {
	Title      string    `json:"title"`
	PriceFrom  float64   `json:"price_from"`
	PriceTo    float64   `json:"price_to"`
	Pagination Pagination `json:"pagination"`
	Category_id string   `json:"category_id"`
	PrType      string   `json:"type"`
}

type ProductList struct {
	Products   []*ProductGet `json:"products"`
	TotalCount int           `json:"total_count"`
	Pagination Pagination    `json:"pagination"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type Product_type struct {
	Prtype string  `json:"prtype"`
	Price  float64 `json:"price"`
}

type Product struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}