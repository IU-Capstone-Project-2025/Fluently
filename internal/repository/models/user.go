package models

type User struct {
	UserID      uint             `gorm:"column:user_id;primaryKey;unique;not null"`
	Name        string           `gorm:"type:varchar(30)"`
	Preferences *UserPreferences `gorm:"foreignKey:UserID;references:UserID"`

	// CreatedAt??
	// UpdatedAt??
}

func (User) TableName() string {
	return "user_account"
}
