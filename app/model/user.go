package model

// User struct
type User struct {
	Base
	UserAPI
}

// UserAPI struct
type UserAPI struct {
	Name     *string `json:"name,omitempty" gorm:"type:varchar(50)"`
	Username *string `json:"username,omitempty" gorm:"type:varchar(50)"`
	Password *string `json:"-" gorm:"type:text"`
}
