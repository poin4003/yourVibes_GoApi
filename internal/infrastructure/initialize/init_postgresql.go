package initialize

import (
	"fmt"
	_ "log"
	"time"

	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"

	"go.uber.org/zap"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/poin4003/yourVibes_GoApi/global"
)

func InitPostgreSql() *gorm.DB {
	m := global.Config.PostgreSql

	dsn := "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
	var s = fmt.Sprintf(dsn, m.Host, m.Port, m.Username, m.Password, m.Dbname, m.SslMode)

	var err error
	db, err := gorm.Open(postgres.Open(s), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	checkErrorPanic(err, "InitPostgreSql initialization error")

	global.Logger.Info("Initializing PostgreSQL Successfully")

	sqlDB, err := db.DB()
	if err != nil {
		global.Logger.Error("InitPostgreSql DB error", zap.Error(err))
	}

	sqlDB.SetConnMaxIdleTime(time.Duration(m.MaxIdleConns) * time.Second)
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(m.ConnMaxLifetime) * time.Second)

	return db
}

func checkErrorPanic(err error, errString string) {
	if err != nil {
		global.Logger.Error(errString, zap.Error(err))
		panic(err)
	}
}

func DBMigrator(db *gorm.DB) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"unaccent\"")
	err := db.AutoMigrate(
		&models.User{},
		&models.Notification{},
		&models.Post{},
		&models.Media{},
		&models.Setting{},
		&models.LikeUserPost{},
		&models.Comment{},
		&models.LikeUserComment{},
		&models.FriendRequest{},
		&models.Friend{},
		&models.NewFeed{},
		&models.Advertise{},
		&models.Bill{},
		&models.Admin{},
		&models.Report{},
		&models.UserReport{},
		&models.PostReport{},
		&models.CommentReport{},
		&models.Conversation{},
		&models.Message{},
		&models.ConversationDetail{},
		&models.Statistics{},
	)
	return err
}
