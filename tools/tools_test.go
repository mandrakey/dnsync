package tools

import (
    "testing"
)

func TestStringInSlice(t *testing.T) {
    sl := []string{"InSlice"}

    if !StringInSlice("InSlice", sl) {
        t.Fatalf("InSlice should be in the slice")
    }
    if StringInSlice("NotInSclice", sl) {
        t.Fatalf("NotInSlice should not be in the slice")
    }
}
