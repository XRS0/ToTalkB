package pkg

type Chat struct {
	Id      string `json:"id" db:"id"`
	Name    string `json:"name" db:"name" binding:"required"`
	OwnerId string `json:"-" db:"created_by"`
}
