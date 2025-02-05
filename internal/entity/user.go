package entity

type User struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password,omitempty"` // Use `omitempty` to avoid exposing passwords in responses
	PhoneNumber string `json:"phone_number"`
	UserRole    string `json:"user_role"`
	AccessToken string `json:"access_token,omitempty"` // Optional field for tokens
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// UserSingleRequest represents a request to fetch a single user.
type UserSingleRequest struct {
	ID       string `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

// UserList represents a paginated list of users.
type UserList struct {
	Items []User `json:"users"`
	Count int    `json:"count"`
}


type UserUpdate struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Password    string    `json:"password,omitempty"` // Ommaviy ko‘rinmasligi uchun
}


