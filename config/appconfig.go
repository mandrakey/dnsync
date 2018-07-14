package config

import (
    "os"
    "fmt"
    "encoding/json"
)

type AppConfig struct {
    ConfigFile string
    Remotes []string
    Verbose bool
    Simulation bool
    Port int
    Host string
    Handlers []Handler
}

type Handler struct {
	Type string
	BindHandler
}

type BindHandler struct {
	BindConfigDir string `json:"config-dir"`
	BindZonefilesPath string `json:"zonefiles-path"`
}

var instance *AppConfig

func AppConfigInstance() *AppConfig {
    if instance == nil {
        instance = &AppConfig{}
    }
    return instance
}

func (ac *AppConfig) LoadFromFile(file string) error {
    if _, err := os.Stat(file); os.IsNotExist(err) {
        return fmt.Errorf("File %s does not exist", file)
    }

    f, err := os.Open(file); if err != nil {
        return fmt.Errorf("Failed to open file: %s", err)
    }

    decoder := json.NewDecoder(f)
    err = decoder.Decode(ac); if err != nil {
        return fmt.Errorf("Failed to decode file: %s", err)
    }

    return nil
}
