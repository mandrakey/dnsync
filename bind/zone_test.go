package bind

import (
    "testing"
)

func TestZoneIsValid(t *testing.T) {
    z := Zone{}
    if z.IsValid() == true {
        t.Fatal("empty zone should be invalid")
    }

    z = Zone{Name: "mydomain.tld."}
    if z.IsValid() == true {
        t.Fatal("zone with only a name should be invalid")
    }

    z = Zone{Masters: []string{"1.2.3.4"}}
    if z.IsValid() == true {
        t.Fatal("zone with only masters should be invalid")
    }

    z = Zone{File: "somefile"}
    if z.IsValid() == true {
        t.Fatal("zone with only file should be invalid")
    }

    z = Zone{Name: "domain.tld.", Masters: []string{"1.2.3.4"}}
    if z.IsValid() == true {
        t.Fatal("zone with only name and masters should be invalid")
    }

    z = Zone{Name: "domain.tld", File: "somefile"}
    if z.IsValid() == true {
        t.Fatal("zone with only name and file should be invalid")
    }

    z = Zone{Masters: []string{"1.2.3.4"}, File: "somefile"}
    if z.IsValid() == true {
        t.Fatal("zone with only masters and file should be invalid")
    }

    z = Zone{Name: "domain.tld", Masters: []string{"1.2.3.4"}, File: "somefile"}
    if z.IsValid() == false {
        t.Fatal("zone with all data should be valid")
    }
}

func TestZoneEquals(t *testing.T) {
    z := &Zone{Name: "domain.tld", Masters: []string{"1.2.3.4", "5.6.7.8"}, File: "somefile"}

    if !z.Equals(z) {
        t.Fatal("the same zone struct does not equal itself")
    }

    if z.Equals(&Zone{Name: "domain2.tld", Masters: []string{"1.2.3.4", "5.6.7.8"}, File: "somefile"}) {
        t.Fatal("different zone names but zones are equal")
    }

    if z.Equals(&Zone{Name: "domain.tld", Masters: []string{"1.2.3.4"}, File: "somefile"}) {
        t.Fatal("different zone masters but zones are equal")
    }

    if z.Equals(&Zone{Name: "domain.tld", Masters: []string{"1.2.3.4", "5.6.7.8"}, File: "someotherfile"}) {
        t.Fatal("different zone files but zones are equal")
    }
}
