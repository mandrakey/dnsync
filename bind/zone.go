package bind

import (
    "fmt"

    "mandrakey.cc/dnsync/tools"
)

type Zone struct {
    Name string
    Masters []string
    File string
}

func CopyZone(zone *Zone) *Zone {
    return &Zone{Name: zone.Name, Masters: zone.Masters, File: zone.File}
}

func (z *Zone) IsValid() bool {
    return !(z.Name == "" || len(z.Masters) == 0 || z.File == "")
}

func (z *Zone) Equals(other *Zone) bool {
    if z.Name != other.Name || z.File != other.File {
        return false
    }

    for _, m := range z.Masters {
        if !tools.StringInSlice(m, other.Masters) {
            return false
        }
    }

    return true
}

func (z *Zone) String() string {
    return fmt.Sprintf("zone{Name: '%s', Masters: %s, File: %s}", z.Name, z.Masters, z.File)
}
