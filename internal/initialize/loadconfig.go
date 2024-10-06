package initialize

import (
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper := viper.New()

	viper.AddConfigPath("./config/")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read configuration %w", err))
	}

	if err := viper.Unmarshal(&global.Config); err != nil {
		fmt.Printf("unable to decode configuration %v", err)
	}
}
