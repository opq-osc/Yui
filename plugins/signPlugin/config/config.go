package config

import (
	"os"
)

//easyjson:json
type Config struct {
	Users  []int64 `json:"users"`
	BotUin int64   `json:"botUin"`
}

var C *Config

func ReadConfig() error {
	fBytes, err := os.ReadFile("/config.json")
	if err != nil {
		return err
	}
	C = &Config{}
	err = C.UnmarshalJSON(fBytes)
	if err != nil {
		return err
	}
	return nil
}
func SaveConfig() error {
	jsonByes, err := C.MarshalJSON()
	if err != nil {
		return err
	}
	err = os.WriteFile("/config.json", []byte(jsonByes), 0777)
	if err != nil {
		return err
	}
	return nil
}
