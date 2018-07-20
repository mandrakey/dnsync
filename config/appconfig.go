package config

import (
    "os"
    "fmt"
    "strings"
    "encoding/json"
)

type AppConfig struct {
    ConfigFile string
    Remotes []string
    Verbose bool
    Logfile string
    Simulation bool
    Port int
    Host string
    Handlers []Handler
}

type Handler struct {
    Name string
	Type string
    BindHandler
}

type BindHandler struct {
	BindConfigFile string `json:"config-file"`
	BindZonefilesPath string `json:"zonefiles-path"`
}

func (h *Handler) UnmarshalJSON(rawdata []byte) error {
    data := make(map[string]string)
    err := json.Unmarshal(rawdata, &data); if err != nil {
        return err
    }

    h.Name = data["name"]
    h.Type = data["type"]

    // BindHandler stuff
    v, ok := data["config-file"]; if ok {
        h.BindConfigFile = strings.TrimSuffix(v, "/")
    }
    v, ok = data["zonefiles-path"]; if ok {
        h.BindZonefilesPath = strings.TrimSuffix(v, "/")
    }

    return nil
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

func (ac *AppConfig) String() string {
    res, err := json.Marshal(ac); if err != nil {
        return fmt.Sprintf("Failed to marshal AppConfig to JSON: %s", err)
    }
    return string(res)
}
