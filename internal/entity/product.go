package entity

type ProductCreate struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SalePrice   float64 `json:"sale_price"`
	Color       string  `json:"color"`
	Size        string  `json:"size"`
	PictureUrl  string  `json:"picture_url"`
	CategoryID string  `json:"category_id"`
}

type ProductCreateForSwagger struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	PictureUrl  string  `json:"picture_url"`
	CategoryID string  `json:"category_id"`
	SalePrice   float64 `json:"sale_price"`
	Color       string  `json:"color"`
	Size        string  `json:"size"`
}

type ProductUpt struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SalePrice   float64 `json:"sale_price"`
	Color       string  `json:"color"`
	Size        string  `json:"size"`
}

type ProductPicture struct {
	ProductId  string `json:"product_id"`
	PictureUrl string `json:"picture_url"`
}

type ProductGet struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SalePrice   float64 `json:"sale_price"`
	Color       string  `json:"color"`
	Size        string  `json:"size"`
	PictureUrls []string`json:"picture_urls"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ProductFilter struct {
	Title       string     `json:"title"`
	PriceFrom   float64    `json:"price_from"`
	PriceTo     float64    `json:"price_to"`
	Pagination  Pagination `json:"pagination"`
	Category_id string     `json:"category_id"`
	PrType      string     `json:"type"`
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

type Product struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
