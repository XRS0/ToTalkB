package pkg

type Chat struct {
    Id      int    `json:"id" db:"id"`
    Name    string `json:"name" db:"name" binding:"required"`
    OwnerId int    `json:"owner_id" db:"created_by"`
}