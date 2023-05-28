package models
import(
	"gorm.io/gorm"
)

type Books struct{
	ID       uint  `gorm:"primary Key; autoIncrement" json:"id"`
	Author   *string `json:"author"`
	Title    *string `json:"title"`
	Publication *string `json:"publication"`
}

//auto migration
func MigrateBooks(db *gorm.DB) error{
	err := db.AutoMigrate(&Books{})
	return err
}

