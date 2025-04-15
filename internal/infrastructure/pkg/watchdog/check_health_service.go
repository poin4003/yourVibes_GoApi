package watchdog

import (
	"encoding/json"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"log"
	"net/http"
	"sync"
	"time"
)

type HealthCheckResponse struct {
	Status       string                 `json:"status"`
	Message      string                 `json:"message"`
	Dependencies map[string]interface{} `json:"dependencies"`
}

type ModerateServiceHealthChecker struct {
	healthURL     string
	timeout       time.Duration
	cron          *cron.Cron
	isRunning     bool
	runningLock   sync.Mutex
	singletonLock sync.Mutex
}

var (
	instance *ModerateServiceHealthChecker
	once     sync.Once
)

func InitWatchdogHealthChecker(healthURL string, timeoutSeconds int) *ModerateServiceHealthChecker {
	once.Do(func() {
		instance = &ModerateServiceHealthChecker{
			healthURL:   healthURL,
			timeout:     time.Duration(timeoutSeconds) * time.Second,
			cron:        cron.New(),
			isRunning:   false,
			runningLock: sync.Mutex{},
		}
	})
	return instance
}

func (mshc *ModerateServiceHealthChecker) StartCronJob() error {
	mshc.singletonLock.Lock()
	defer mshc.singletonLock.Unlock()

	mshc.runningLock.Lock()
	if mshc.isRunning {
		mshc.runningLock.Unlock()
		global.Logger.Warn("Python health check cronjob is already running")
		return nil
	}
	mshc.isRunning = true
	mshc.runningLock.Unlock()

	_, err := mshc.cron.AddFunc("*/20 * * * *", mshc.checkHealth)
	if err != nil {
		mshc.runningLock.Lock()
		mshc.isRunning = false
		mshc.runningLock.Unlock()
		global.Logger.Error("Failed to add Python health check cronjob", zap.Error(err))
		return err
	}

	mshc.checkHealth()

	mshc.cron.Start()
	global.Logger.Info("Python health check cronjob started", zap.String("schedule", "every 20 minutes"))
	return nil
}

func (mshc *ModerateServiceHealthChecker) StopCronJob() {
	mshc.runningLock.Lock()
	defer mshc.runningLock.Unlock()

	if !mshc.isRunning {
		global.Logger.Warn("Python health check cronjob is not running")
		return
	}

	mshc.cron.Stop()
	mshc.isRunning = false
	global.Logger.Info("Python health check cronjob stopped")
}

func (mshc *ModerateServiceHealthChecker) checkHealth() {
	global.Logger.Info("Checking Python server health", zap.String("url", mshc.healthURL))

	client := &http.Client{
		Timeout: mshc.timeout,
	}

	resp, err := client.Get(mshc.healthURL)
	if err != nil {
		global.Logger.Error("Failed to connect to Python server", zap.Error(err))
		log.Panicf("Failed to connect to Python server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		global.Logger.Error("Python server is unhealthy", zap.Int("status_code", resp.StatusCode))
		log.Panicf("Python server is unhealthy: status code %d", resp.StatusCode)
	}

	var healthResp HealthCheckResponse
	if err = json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		global.Logger.Error("Failed to parse Python server health response", zap.Error(err))
		log.Panicf("Failed to parse Python server health response: %v", err)
	}

	if healthResp.Status != "healthy" {
		global.Logger.Error("Python server reported unhealthy status",
			zap.String("message", healthResp.Message),
			zap.Any("dependencies", healthResp.Dependencies))
		log.Panicf("Python server is unhealthy: %s", healthResp.Message)
	}

	global.Logger.Info("Python server is healthy",
		zap.String("message", healthResp.Message),
		zap.Any("dependencies", healthResp.Dependencies),
	)
}
