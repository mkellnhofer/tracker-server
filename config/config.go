package config

import (
	"log"

	"github.com/go-ini/ini"
)

type Config struct {
	Port int
}

func LoadConfig() *Config {
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		log.Fatal("Config file missing!")
	}

	port := getIntValue(cfg, "server", "port")

	return &Config{port}
}

func getStringValue(file *ini.File, secName string, keyName string) string {
	return getKey(file, secName, keyName).String()
}

func getIntValue(file *ini.File, secName string, keyName string) int {
	val, err := getKey(file, secName, keyName).Int()
	if err != nil {
		log.Fatalf("Config file has invalid value for key '%s'!", keyName)
	}
	return val
}

func getKey(file *ini.File, secName string, keyName string) *ini.Key {
	sec, err := file.GetSection(secName)
	if err != nil {
		log.Fatalf("Config file missing section '%s'!", secName)
	}

	if !sec.HasKey(keyName) {
		log.Fatalf("Config file missing key '%s'!", keyName)
	}

	return sec.Key(keyName)
}
