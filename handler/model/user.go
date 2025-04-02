package model

type User struct {
	Name     string `json:"name" binding:"required"`
	PassWord string `json:"password" binding:"required"`
}

type ModifyPassRequest struct {
	OldPass string `json:"old_pass" binding:"required,len=32"`
	NewPass string `json:"new_pass" binding:"required,len=32"`
}
