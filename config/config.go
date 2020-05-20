package config

import (
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//Init function to initialize all configurations
func Init(path string) error {
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

//GetString method to get configs from config file
func GetString(keyName string) string {
	return viper.GetString(keyName)
}

//GetStringSlice method to get configs from config file
func GetStringSlice(keyName string) []string {
	return viper.GetStringSlice(keyName)
}

func GetInt(keyName string) int {
	return viper.GetInt(keyName)
}
