package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Connection struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Config struct {
	Connections map[string]Connection `yaml:"connections"`
}

func readConfig() (*Config, error) {
	// Open the file
	filePath := filepath.Join(os.Getenv("HOME"), ".config", "lazydb.yml")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error opening config file, make sure it exists at %s", filePath)
	}
	defer file.Close()

	// Decode the file
	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// return array of name and url
func GetConnections() (map[string]Connection, error) {
	config, err := readConfig()
	if err != nil {
		return nil, err
	}

	return config.Connections, nil
}

func (c *Connection) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.User, c.Password, c.Host, c.Port, c.Database)
}
