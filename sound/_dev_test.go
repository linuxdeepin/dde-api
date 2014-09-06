package main

import (
	"testing"
	"time"
)

func TestPlaySound(t *testing.T) {
	s := &Sound{}
	err := s.PlaySystemSound("bell")
	if err != nil {
		t.Error(err)
	}
	time.Sleep(2 * time.Second)
}
