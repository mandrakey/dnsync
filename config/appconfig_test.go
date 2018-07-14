package config

import (
    "testing"
)

func TestLoadFromFile(t *testing.T) {
    ac := AppConfig{}
    err := ac.LoadFromFile("./appconfig_test.json"); if err != nil {
        t.Fatalf("Failed loading config: %s", err)
    }

    if len(ac.Remotes) != 2 {
        t.Fatalf("Wrong amount of remotes loaded")
    }
    if ac.Remotes[0] != "127.0.0.1" {
        t.Fatalf("First remote is not 127.0.0.1")
    }
    if ac.Remotes[1] != "1.2.3.4" {
        t.Fatalf("Second remote is not 1.2.3.4")
    }

    if ac.Port != 53001 {
        t.Fatalf("Port is not 53001")
    }

    if ac.Host != "0.0.0.0" {
        t.Fatalf("Host is not 0.0.0.0")
    }

    if len(ac.Handlers) != 1 {
        t.Fatalf("Amount of handlers is not 1")
    }

    if ac.Handlers[0].Type != "bind" {
        t.Fatalf("First handler type is not bind")
    }
    if ac.Handlers[0].BindConfigDir != "config1" {
        t.Fatalf("First handler bind config-dir is not config1")
    }
    if ac.Handlers[0].BindZonefilesPath != "path1" {
        t.Fatalf("First handler bind zonefiles-path is not path1")
    }
}
