package config

import (
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//InitConfig function to initialize all configurations
func InitConfig(path string) error {

	viper.SetConfigName(path)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		logger.WithField("error config file", err.Error()).Error("Cannot initialize config")
		return err
	}
	logger.WithField("msg", "initialized successfully").Info("Config initialization")
	return nil
}

//GetConfig method to get configs from config file
func GetConfig(keyName string) string {
	keyValue := viper.GetString(keyName)
	return keyValue
}

//GetStringSlice method to get configs from config file
func GetStringSlice(keyName string) []string {
	keyValue := viper.GetStringSlice(keyName)
	return keyValue
}
