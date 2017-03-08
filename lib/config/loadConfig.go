package config

import (
	"../tools"
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	DbObj        tools.Database `json:"Database"`
}

func (c *Config) LoadConfig(configFile string) error {

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, c)
	if err != nil {
		return err
	}

	c.DbObj.Host = c.DbObj.Host

	return nil
}
