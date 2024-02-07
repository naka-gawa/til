package main

import "testing"

func TestHelloy(t *testing.T) {
    got := Helloy("Chris")
    want := "Hello, Chris"

    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
