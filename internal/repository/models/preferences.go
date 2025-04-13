package models

type CEFRLevel string

const (
	A1 CEFRLevel = "A1"
	A2 CEFRLevel = "A2"
	B1 CEFRLevel = "B1"
	B2 CEFRLevel = "B2"
	C1 CEFRLevel = "C1"
	C2 CEFRLevel = "C2"
)

type UserPreferences struct {
	UserID          uint      `gorm:"primaryKey;not null;unique;column:user_id"`
	CEFRLevel       CEFRLevel `gorm:"type:varchar(2);not null"`
	Goal            string    `gorm:"type:varchar(100);default:generally"`
	Notifications   bool      `gorm:"default:true"`
	Advertisments   bool      `gorm:"default:true"`
	WordsPerDay     int       `gorm:"default:10"`
	NotificationsAt string    `gorm:"type:varchar(5);default:'09:00'"`

	User User `gorm:"foreignKey:UserID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (UserPreferences) TableName() string {
	return "user_preferences"
}
