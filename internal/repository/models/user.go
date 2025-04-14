package models

<<<<<<< HEAD
import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name     string    `gorm:"type:varchar(100);not null"`
	SubLevel bool      `gorm:"default:false"`
	PrefID   uuid.UUID `gorm:"type:uuid"`

	Pref Preference `gorm:"foreignKey:PrefID;constraint:OnUpdate:CASCADE,OnDelete: SET NULL"`
=======
type User struct {
	UserID      uint             `gorm:"column:user_id;primaryKey;unique;not null"`
	Name        string           `gorm:"type:varchar(30)"`
	Preferences *UserPreferences `gorm:"foreignKey:UserID;references:UserID"`
>>>>>>> d67dbcc (Add all user logic)
}

func (User) TableName() string {
	return "users"
}
