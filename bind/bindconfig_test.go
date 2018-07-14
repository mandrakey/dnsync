package bind

import (
    "testing"
)

func TestBindConfigLoad(t *testing.T) {
    bc := NewBindConfig()
    bc.Load("./test.conf")
    out := bc.ToString()
    expected := "Zone = { name: 'mjui.de', masters: [ 88.99.47.253 ], file: \"/etc/bind/db.mjui.de\" };\n" +
        "Zone = { name: 'dau.fun', masters: [ 88.99.47.253 ], file: \"/etc/bind/db.dau.fun\" };\n"

    if out != expected {
        t.Fatalf("Zone output not as expected.\nExpect: %s\nActual: %s\n", expected, out)
    }
}

func TestBindConfigSave(t *testing.T) {
    bc := NewBindConfig()
    bc2 := NewBindConfig()
    file1 := "./test.conf"
    file2 := "./test2.conf"

    if !bc.Equals(bc2) {
        t.Fatal("empty bind config instances should be equal")
    }

    bc.Load(file1)
    bc.Save(file2)

    // Load it again and compare
    bc2.Load(file2)

    if !bc.Equals(bc2) {
        t.Fatal("saved and re-loaded bind config not equal to original")
    }
}

func TestBindConfigAddZone(t *testing.T) {
    bc := NewBindConfig()
    bc2 := NewBindConfig()
    z1 := &Zone{Name: "domain.tld", Masters: []string{"1.2.3.4"}, File: "somefile"}
    //z2 := &Zone{Name: "domain2.tld", Masters: []string{"1.2.3.4"}, File: "somefile"}

    bc.AddZone(z1)
    if bc.Equals(bc2) {
        t.Fatal("bind confings are equal after adding a zone to only one")
    }

    bc2.AddZone(z1)
    if !bc.Equals(bc2) {
        t.Fatal("bind configs are not equal after adding the same zone to the second config")
    }
}

func TestBindConfigRemoveZone(t *testing.T) {
    bc := NewBindConfig()
    bc2 := NewBindConfig()
    z1 := &Zone{Name: "domain.tld", Masters: []string{"1.2.3.4"}, File: "somefile"}

    bc.AddZone(z1)
    if bc.Equals(bc2) {
        t.Fatal("bind confings are equal after adding a zone to only one")
    }

    bc.RemoveZone(z1)
    if !bc.Equals(bc2) {
        t.Fatal("bind configs are not equal after removing a zone again")
    }
}

func TestBindConfigGetZone(t *testing.T) {
    bc := NewBindConfig()
    z1 := &Zone{Name: "domain.tld", Masters: []string{"1.2.3.4"}, File: "somefile"}
    bc.AddZone(z1)

    z := bc.GetZone(z1.Name)
    if !z1.Equals(z) {
        t.Fatal("retrieving previously added zone from config yields different zone")
    }
}

func TestBindConfigEquals(t *testing.T) {
    bc := NewBindConfig()
    bc2 := NewBindConfig()

    bc.AddZone(&Zone{Name: "domain.tld", Masters: []string{"1.2.3.4"}, File: "somefile"})
    bc.AddZone(&Zone{Name: "domain2.tld", Masters: []string{"1.2.3.4"}, File: "somefile"})
    bc2.AddZone(&Zone{Name: "domain2.tld", Masters: []string{"1.2.3.4"}, File: "somefile"})
    bc2.AddZone(&Zone{Name: "domain.tld", Masters: []string{"1.2.3.4"}, File: "somefile"})

    if !bc.Equals(bc2) {
        t.Fatal("BindConfig instances based off the same file are not equal")
    }
}
