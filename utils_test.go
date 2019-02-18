package main

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestReg(t *testing.T) {

}

func TestUA(t *testing.T) {
    ua1 := RandomUserAgent()
    ua2 := RandomUserAgent()
    t.Logf("ua1: %s, ua2: %s", ua1, ua2)
    assert.NotEqual(t, ua1, ua2)
}
