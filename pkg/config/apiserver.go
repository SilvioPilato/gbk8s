package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
    Http struct {
        Port int `yaml:"port"`
        Host string `yaml:"host"`
    } `yaml:"http"`
    Agent struct {
        Host string `yaml:"Host"`
        Port int `yaml:"port"`
    } `yaml:"agent"`
}


func ReadYamlConfig(cfg *Config, filePath string) (error) {
    f, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer f.Close()

    decoder := yaml.NewDecoder(f)
    err = decoder.Decode(cfg)
    if err != nil {
		return err
    }
	return nil
} 