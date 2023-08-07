package db

import (
	"clipboard/models"
	"clipboard/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"time"
)

var db *gorm.DB

func GetDb() *gorm.DB {
	return db
}

func CheckDbExist(filepath string) (bool, error) {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func CreateDb(filepath string) (*os.File, error) {
	return os.Create(filepath)
}

func OpenDb(filepath string) (*gorm.DB, error) {
	var err error
	db, err = gorm.Open(sqlite.Open(filepath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitDb(filepath string) (db *gorm.DB, err error) {
	db, err = OpenDb(filepath)
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&models.ClipBoardEntity{}); err != nil {
		return nil, err
	}
	return db, nil
}

func AddClipBoard(content string) error {
	hash := utils.GetHash(content)
	entity := models.ClipBoardEntity{
		Content: content,
		Hash:    hash,
		Time:    time.Now().String(),
	}
	return db.Create(&entity).Error
}

func GetClipBoard(offset, limit int) ([]models.ClipBoardEntity, error) {
	var entities []models.ClipBoardEntity
	return entities, db.Offset(offset).Limit(limit).Find(&entities).Error
}

func GetLatestClipBoard() (models.ClipBoardEntity, error) {
	var entity models.ClipBoardEntity
	return entity, db.Last(&entity).Error
}
