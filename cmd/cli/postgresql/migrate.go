package main

import (
	"github.com/poin4003/yourVibes_GoApi/global"
	initialize2 "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/initialize"
	"go.uber.org/zap"
	"log"
)

func main() {
	initialize2.LoadConfig()
	initialize2.InitLogger()
	initialize2.InitPostgreSql()

	logger := global.Logger

	logger.Info("Starting migration process...")
	if err := initialize2.DBMigrator(global.Pdb); err != nil {
		logger.Error("Unable to migrate database", zap.Error(err))
		log.Fatalln("Migration failed:", err)
	} else {
		logger.Info("Migration complete successfully")
	}

	logger.Info("Migration process finished.")
}
