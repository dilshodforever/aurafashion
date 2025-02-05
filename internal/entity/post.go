package entity

type PostCreate struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	PictureUrls string `json:"picture_urls"`
}

type PostGet struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	PictureUrls []string `json:"picture_urls"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string
}



type PostUpdate struct {
	ID   string  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PostList struct {
	Posts      []PostGet `json:"posts"`
	TotalCount int     `json:"total_count"`
	Pagination Pagination `json:"pagination"`
}

type PostPicture struct {
	PostID    string `json:"post_id"`
	PictureUrl string `json:"picture_url"`
}

type PostFilter struct {
	Title       string `json:"title"`
	CreatedFrom string `json:"created_from"`
	CreatedTo   string `json:"created_to"`
	Limit       int    `json:"limit"`
	Page        int    `json:"page"`
}

