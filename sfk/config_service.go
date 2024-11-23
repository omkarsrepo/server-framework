package sfk

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

var (
	configServiceInstance *configService
	configServiceOnce     sync.Once
	configFiles           = map[string]string{
		"default": "../config/default.json",
		"sandbox": "../config/sandbox.json",
		"prod":    "../config/prod.json",
	}
	configFileType = "json"
)

type configService struct {
	viper *viper.Viper
}

type ConfigService interface {
	GetString(key string) string
	GetInt(key string) int
	GetInt64(key string) int64
	GetBool(key string) bool
	GetViper() *viper.Viper
}

func buildConfig(isDefault bool) *viper.Viper {
	config := viper.New()
	config.SetConfigType(configFileType)
	env := "default"

	if isDefault {
		env = "default"
	} else {
		env = viper.GetString("env")
	}

	configFilePath, found := configFiles[env]
	if !found {
		panic(fmt.Sprintf(`Config file for %s env not found`, env))
	}

	config.SetConfigFile(configFilePath)

	if err := config.ReadInConfig(); err != nil {
		panic(fmt.Sprintf(`Error reading config file for %s environment. Error %s`, env, err))
	}

	return config
}

func getViper() *viper.Viper {
	defaultConfig := buildConfig(true)
	envConfig := buildConfig(false)

	err := viper.MergeConfigMap(defaultConfig.AllSettings())
	if err != nil {
		panic(err)
	}

	err = viper.MergeConfigMap(envConfig.AllSettings())
	if err != nil {
		panic(err)
	}

	return viper.GetViper()
}

func ConfigServiceInstance() ConfigService {
	configServiceOnce.Do(func() {
		configServiceInstance = &configService{
			viper: getViper(),
		}
	})

	return configServiceInstance
}

func (c *configService) GetString(key string) string {
	return c.viper.GetString(key)
}

func (c *configService) GetInt(key string) int {
	return c.viper.GetInt(key)
}

func (c *configService) GetInt64(key string) int64 {
	return c.viper.GetInt64(key)
}

func (c *configService) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

func (c *configService) GetViper() *viper.Viper {
	return c.viper
}
