package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"queding.com/go/common/config"
)

var (
	db *gorm.DB
)

func init() {
	url := config.GetString("mysql.connect")

	var err error
	db, err = gorm.Open("mysql", url)
	if err != nil {
		panic("failed to connect database")
	}
}

func Close() {
	db.Close()
}
