/* This file is part of DNSync.
 *
 * Copyright (C) 2018 Maurice Bleuel <mandrakey@bleuelmedia.com>
 * Licensed undert the simplified BSD license. For further details see COPYING.
 */

package bind

import (
    "os"
    "fmt"
    "bufio"
    "regexp"
    "strings"
)

// Represents a bind zone config file containing one or more zones.
type BindConfig struct {
    zones map[string]*Zone
}

var rxZone *regexp.Regexp
var rxEnd *regexp.Regexp
var rxMastersBegin *regexp.Regexp
var rxFile *regexp.Regexp
var rxMaster *regexp.Regexp

// Creates a new empty BindConfig instance and returns a pointer to it. Note: This function should be used to
// generate BindConfig instances, as it also prepares the regular expressions necessary for parsing.
// todo: Find a better place for compiling the regular expressions ...
func NewBindConfig() *BindConfig {
    if rxZone == nil {
        rxZone = regexp.MustCompile("^zone \"(.+)\" \\{")
        rxEnd = regexp.MustCompile("\\};")
        rxMastersBegin = regexp.MustCompile("masters \\{")
        rxFile = regexp.MustCompile("file \"(.+)\";")
        rxMaster = regexp.MustCompile("(\\S+);")
    }

    return &BindConfig{zones: make(map[string]*Zone)}
}

// Load a BindConfig from a given bind configuration file and store it in the current instance.
func (bc *BindConfig) Load(file string) error {
    if _, err := os.Stat(file); os.IsNotExist(err) {
        return fmt.Errorf("The given file %s does not exist.\n", file)
    }

    bc.zones = make(map[string]*Zone)
    f, err := os.Open(file); if err != nil {
        return fmt.Errorf("Failed to open file: %s\n", err)
    }

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()

        m := rxZone.FindStringSubmatch(line)
        if len(m) == 0 {
            continue
        }

        zone := bc.parseZone(scanner, m[1])
        bc.zones[zone.Name] = zone
    }
    return nil
}

// Save the current BindConfig instance into a specified file to become a bind configuration file. Already existing
// files will be overwritten.
func (bc *BindConfig) Save(file string) error {
    f, err := os.Create(file); if err != nil {
        return fmt.Errorf("Failed to open file: %s\n", err)
    }

    for _, zone := range bc.zones {
        f.WriteString(fmt.Sprintf("zone \"%s\" {\n", zone.Name))
        f.WriteString("        type slave;\n")
        f.WriteString("        masters {\n")

        for _, m := range zone.Masters {
            f.WriteString(fmt.Sprintf("                %s;\n", m))
        }

        f.WriteString("                };\n")
        f.WriteString(fmt.Sprintf("        file \"%s\";\n", zone.File))
        f.WriteString("};\n")
    }

    return nil
}

// Adds a given zone to the current BindConfig. Already existing zones will be replaced.
func (bc *BindConfig) AddZone(zone *Zone) {
    bc.zones[zone.Name] = zone
}

// Remove a given zone from the current BindConfig, if it contains the zone.
func (bc *BindConfig) RemoveZone(zone *Zone) {
    delete(bc.zones, zone.Name)
}

// Retrieve the zone instance for a given domain name from this BindConfig.
func (bc *BindConfig) GetZone(name string) *Zone {
    o := bc.zones[name]; if o == nil {
        return nil
    }
    return CopyZone(o)
}

// Create a string representation of this BindConfig instance.
func (bc *BindConfig) String() string {
    res := make([]string, 0)

    for _, zone := range bc.zones {
        res = append(
            res,
            fmt.Sprintf(
                "Zone = { name: '%s', masters: [ %s ], file: \"%s\" };\n",
                zone.Name,
                strings.Join(zone.Masters, ", "),
                zone.File,
            ),
        )
    }

    return strings.Join(res, "")
}

// Parses given data from scanner and stores the information in this BindConfig.
func (bc *BindConfig) parseZone(scanner *bufio.Scanner, name string) *Zone {
    z := Zone{Name: name}

    inMasters := false
    for scanner.Scan() {
        line := scanner.Text()

        if rxEnd.MatchString(line) {
            if inMasters {
                inMasters = false
            } else {
                return &z
            }
        }

        if !inMasters && rxMastersBegin.MatchString(line) {
            inMasters = true
            continue
        }

        if !inMasters && rxFile.MatchString(line) {
            m := rxFile.FindStringSubmatch(line)
            z.File = strings.TrimSpace(m[1])
        }

        if inMasters {
            m := rxMaster.FindStringSubmatch(line)
            if len(m) == 0 {
                continue
            }
            z.Masters = append(z.Masters, strings.TrimSpace(m[1]))
        }
    }

    return &z
}

// Check whether or not this BindConfig is the same as other.
func (bc *BindConfig) Equals(other *BindConfig) bool {
    matches := make(map[string]bool)

    outer:
    for _, z := range bc.zones {
        matches[z.Name] = false
        for _, z2 := range other.zones {
            if !z.Equals(z2) {
                continue
            }
            matches[z.Name] = true
            continue outer
        }
    }

    for _, matched := range matches {
        if !matched {
            return false
        }
    }
    return true
}
