/* This file is part of DNSync.
 *
 * Copyright (C) 2018 Maurice Bleuel <mandrakey@bleuelmedia.com>
 * Licensed undert the simplified BSD license. For further details see COPYING.
 */

package bind

import (
    "fmt"

    "mandrakey.cc/dnsync/tools"
)

// Represents a bind domain zone.
type Zone struct {
    Name string
    Masters []string
    File string
}

// Create a new Zone instance based on zone.
func CopyZone(zone *Zone) *Zone {
    return &Zone{Name: zone.Name, Masters: zone.Masters, File: zone.File}
}

// Check whether this Zone instance contains all necessary information to be a valid, working DNS zone.
func (z *Zone) IsValid() bool {
    return !(z.Name == "" || len(z.Masters) == 0 || z.File == "")
}

// Check whether or not this Zone contains the same information as other.
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

// Create a string representation of this Zone.
func (z *Zone) String() string {
    return fmt.Sprintf("zone{Name: '%s', Masters: %s, File: %s}", z.Name, z.Masters, z.File)
}
