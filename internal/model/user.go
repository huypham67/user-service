package model

// User represents the user model in the database.
type User struct {
	BaseModel

	DisplayName string `json:"display_name" gorm:"not null;column:display_name"`
	Username    string `json:"username" gorm:"not null;unique;column:username"`
	Email       string `json:"email" gorm:"unique;not null;column:email"`
	Password    string `json:"-" gorm:"not null;column:password"`
}
