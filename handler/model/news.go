package model

type News struct {
	Title   string `json:"title"  binding:"required,gte=1" validate:"required,gte=1"`   //长度>=1
	Content string `json:"content"  binding:"required,gte=1" validate:"required,gte=1"` //长度>=1
	Id      int    `json:"id"`
}
