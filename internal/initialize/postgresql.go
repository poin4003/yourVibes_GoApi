package initialize

import (
	"fmt"
	_ "log"
	"time"

	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"go.uber.org/zap"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/poin4003/yourVibes_GoApi/global"
)

func InitPostgreSql() {
	m := global.Config.PostgreSql

	dsn := "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
	var s = fmt.Sprintf(dsn, m.Host, m.Port, m.Username, m.Password, m.Dbname, m.SslMode)

	var err error
	db, err := gorm.Open(postgres.Open(s), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	checkErrorPanic(err, "InitPostgreSql initialization error")

	global.Pdb = db
	global.Logger.Info("Initializing PostgreSQL Successfully")

	SetPool()
}

func checkErrorPanic(err error, errString string) {
	if err != nil {
		global.Logger.Error(errString, zap.Error(err))
		panic(err)
	}
}

func SetPool() {
	m := global.Config.PostgreSql

	sqlDb, err := global.Pdb.DB()
	checkErrorPanic(err, "Failed to get PostgreSql")

	sqlDb.SetConnMaxIdleTime(time.Duration(m.MaxIdleConns) * time.Second)
	sqlDb.SetMaxOpenConns(m.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(time.Duration(m.ConnMaxLifetime) * time.Second)
}

func DBMigrator(db *gorm.DB) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"unaccent\"")
	err := db.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Media{},
		&model.Setting{},
		&model.LikeUserPost{},
	)
	return err
}
