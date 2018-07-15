package bind

import (
    "os"
    "fmt"
    "bufio"
    "regexp"
    "strings"
)

type BindConfig struct {
    zones map[string]*Zone
}

var rxZone *regexp.Regexp
var rxEnd *regexp.Regexp
var rxMastersBegin *regexp.Regexp
var rxFile *regexp.Regexp
var rxMaster *regexp.Regexp

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

func (bc *BindConfig) AddZone(zone *Zone) {
    bc.zones[zone.Name] = zone
}

func (bc *BindConfig) RemoveZone(zone *Zone) {
    delete(bc.zones, zone.Name)
}

func (bc *BindConfig) GetZone(name string) *Zone {
    o := bc.zones[name]; if o == nil {
        return nil
    }
    return CopyZone(o)
}

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
