package entity

type CategoryRes struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CategoryId struct {
	ID string `json:"id"`
}

type CategoryListsReq struct {
	Name   string     `json:"name"`
	Filter Pagination `json:"filter"`
}

type CategoryListsRes struct {
	Categories  []CategoryRes `json:"Categorys"`
	TotalCount int32         `json:"total_count"`
}

type CategoryUptBody struct {
	Name        string  `json:"name"`
}

type CategoryUpt struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
}
