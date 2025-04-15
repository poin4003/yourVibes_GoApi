package initialize

import (
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/watchdog"
	"log"
)

func InitWatchDog() {
	healthChecker := watchdog.InitWatchdogHealthChecker(global.Config.ModerateService.HealthURL, 5)
	if err := healthChecker.StartCronJob(); err != nil {
		log.Panicf("Failed to start Python health check cronjob: %v", err)
	}
	defer healthChecker.StopCronJob()
}
