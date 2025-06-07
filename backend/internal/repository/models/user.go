package models

type User struct {
	UserID      uint        `gorm:"column:user_id;primaryKey;unique;not null"`
	Name        string      `gorm:"type:varchar(30)"`
	Preferences *Preference `gorm:"foreignKey:UserID;references:UserID"`
}

func (User) TableName() string {
	return "users"
}
