package database

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("database/demo.db"))
	if err != nil {
		log.Fatal().Msg("database open error")
	}
	if DB.Error != nil {
		log.Fatal().Msg("database error")
	}
}
