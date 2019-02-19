package main

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "math/rand"
    "regexp"
    "testing"
)

func TestReg(t *testing.T) {
    var reg = regexp.MustCompile("{uint}")
    var o = "dalfjoi/{uint}"
    var n = reg.ReplaceAllString(o, fmt.Sprint(rand.Intn(64)))
    fmt.Printf("%s => %s\n", o, n)
    n = reg.ReplaceAllString(o, fmt.Sprint(rand.Intn(64)))
    fmt.Printf("%s => %s\n", o, n)
    o = n
    n = reg.ReplaceAllString(o, fmt.Sprint(rand.Intn(64)))
    fmt.Printf("%s => %s\n", o, n)
    o = "dalfjoi/"
    n = reg.ReplaceAllString(o, fmt.Sprint(rand.Intn(64)))
    fmt.Printf("%s => %s\n", o, n)
}

func TestUA(t *testing.T) {
    ua1 := RandomUserAgent()
    ua2 := RandomUserAgent()
    t.Logf("ua1: %s, ua2: %s", ua1, ua2)
    assert.NotEqual(t, ua1, ua2)
}
