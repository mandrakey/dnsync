/* This file is part of DNSync.
 *
 * Copyright (C) 2018 Maurice Bleuel <mandrakey@bleuelmedia.com>
 * Licensed undert the simplified BSD license. For further details see COPYING.
 */

package config

import (
    "os"
    "fmt"
    "strings"
    "encoding/json"
)

// Represents the application configuration.
type AppConfig struct {
    ConfigFile string
    Remotes []string
    Verbose bool
    Logfile string
    Loglevel string
    Simulation bool
    Port int
    Host string
    Handlers []Handler
}

// Basic DNS server handler struct containing BindHandler fields.
type Handler struct {
    Name string
	Type string
    BindHandler
}

// Special fields struct for bind server handlers.
type BindHandler struct {
	BindConfigFile string `json:"config-file"`
	BindZonefilesPath string `json:"zonefiles-path"`
}

// Unmarshal Handler JSON data read from a configuration file.
// Mainly used to extract and populate specialised handler fields.
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

// The global instance of the AppConfig struct.
var instance *AppConfig

// Get the global AppConfig instance. If it does not exist yet, it will be created.
func AppConfigInstance() *AppConfig {
    if instance == nil {
        instance = &AppConfig{Loglevel: "info"}
    }
    return instance
}

// Populate the fields of this AppConfig by reading data from a given file. The file must be JSON.
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

// Retrieve a JSON formatted string representation of the current AppConfig.
func (ac *AppConfig) String() string {
    res, err := json.Marshal(ac); if err != nil {
        return fmt.Sprintf("Failed to marshal AppConfig to JSON: %s", err)
    }
    return string(res)
}
