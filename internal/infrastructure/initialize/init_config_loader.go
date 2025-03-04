package initialize

import (
	"fmt"

	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper := viper.New()
	viper.AutomaticEnv()

	viper.AddConfigPath("./config/")
	//viper.SetConfigName(viper.GetString("YOURVIBES_SERVER_CONFIG_FILE"))
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read configuration %w", err))
	}

	if err := viper.Unmarshal(&global.Config); err != nil {
		fmt.Printf("unable to decode configuration %v", err)
	}
}
