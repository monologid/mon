package mon

import (
	"fmt"

	"github.com/spf13/viper"
)

// NewConfig sets config using yaml format
func NewConfig(folderPath string, config interface{}) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if len(folderPath) != 0 {
		viper.AddConfigPath(folderPath)
	}

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = viper.Unmarshal(config)
	if err != nil {
		panic(fmt.Errorf("Failed to parse config file: %s \n", err))
	}
}
