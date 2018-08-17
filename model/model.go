package model

type Post struct {
	ID       uint      `gorm:"primary_key" json:"id"`
	Title    string    `gorm:"not null" json:"title" validate:"required"`
	Body     string    `gorm:"not null" json:"body" validate:"required"`
	Comments []Comment `json:"comments"`
}

type Comment struct {
	ID     uint   `gorm:"primary_key" json:"id"`
	Name   string `gorm:"not null" json:"name" validate:"required"`
	Body   string `gorm:"not null" json:"body" validate:"required"`
	PostID uint   `gorm:"type:int REFERENCES posts(id) ON DELETE CASCADE" json:"postId"`
}
