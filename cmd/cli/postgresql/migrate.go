package main

import (
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/initialize"
	"go.uber.org/zap"
	"log"
)

func main() {
	initialize.LoadConfig()
	initialize.InitLogger()
	initialize.InitPostgreSql()

	logger := global.Logger

	logger.Info("Starting migration process...")
	if err := initialize.DBMigrator(global.Pdb); err != nil {
		logger.Error("Unable to migrate database", zap.Error(err))
		log.Fatalln("Migration failed:", err)
	} else {
		logger.Info("Migration complete successfully")
	}

	logger.Info("Migration process finished.")
}
