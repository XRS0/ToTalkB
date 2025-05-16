package pkg

type User struct {
	Id       string `json:"-" db:"id"`
	Login    string `json:"login" db:"login" binding:"required"`
	Password string `json:"password" db:"password_hash" binding:"required"`
	Name     string `json:"name" db:"name" binding:"required"`
	Role     string `json:"role" db:"role" binding:"required"`
}

type UserResponse struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type SignInRequest struct {
	Login    string `json:"login" db:"login" binding:"required"`
	Password string `json:"password" db:"password_hash" binding:"required"`
}
