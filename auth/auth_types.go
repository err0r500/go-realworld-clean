package auth

// User represents a RealWorld user
type User struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

type UserResponse struct {
	User User `json:"user"`
}
