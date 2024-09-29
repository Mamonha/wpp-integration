package services

import (
	"fmt"
	"wpp-integration/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(dsn string) error {
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("erro ao conectar ao banco de dados: %v", err)
	}

	if err := db.AutoMigrate(&models.MessageRecord{}); err != nil {
		return fmt.Errorf("erro ao migrar o schema: %v", err)
	}

	return nil
}

func SaveMessageToDB(message models.MessageRecord) error {
	fmt.Println(message)
	if err := db.Create(&message).Error; err != nil {
		return fmt.Errorf("erro ao salvar mensagem no banco de dados: %v", err)
	}
	return nil
}

func GetDB() *gorm.DB {
	return db
}
