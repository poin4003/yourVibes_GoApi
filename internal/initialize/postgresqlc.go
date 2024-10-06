package initialize

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/poin4003/yourVibes_GoApi/global"
	"go.uber.org/zap"
	"time"
)

func InitPostgreSqlc() {
	m := global.Config.PostgreSql

	dsn := "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
	var s = fmt.Sprintf(dsn, m.Host, m.Port, m.Username, m.Password, m.Dbname, m.SslMode)

	db, err := sql.Open("postgres", s)

	checkErrorPanic(err, "InitPostgreSqlc initialization error")

	global.Logger.Info("Initializing MySQL Successfully")

	global.Pdb = db

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

	sqlDb := global.Pdb

	sqlDb.SetConnMaxIdleTime(time.Duration(m.MaxIdleConns))
	sqlDb.SetMaxOpenConns(m.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(time.Duration(m.ConnMaxLifetime))
}
