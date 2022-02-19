package hfsm

import (
	"testing"
	"time"
)

func TestNewState(t *testing.T) {
	state := NewState("test")
	if state.GetName() != "test" {
		t.Error("state name not be test")
	}
}

func TestWithOnEnter(t *testing.T) {
	cb := make(chan struct{}, 1)
	state := NewState("test", WithOnEnter(func(pre string, ctx interface{}) {
		cb <- struct{}{}
	}))
	state.enter("", nil)
	select {
	case <-cb:
	case <-time.After(time.Duration(1)):
		t.Error("failed to raise onEnter")
	}
}

func TestWithOnUpdate(t *testing.T) {
	cb := make(chan struct{}, 1)
	state := NewState("test", WithOnUpdate(func(deltaTime time.Duration, ctx interface{}) {
		cb <- struct{}{}
	}))
	state.update(time.Duration(1), nil)
	select {
	case <-cb:
	case <-time.After(time.Duration(1)):
		t.Error("failed to raise onupdate")
	}
}

func TestWithOnExit(t *testing.T) {
	cb := make(chan struct{}, 1)
	state := NewState("test", WithOnExit(func(pre string, ctx interface{}) {
		cb <- struct{}{}
	}))
	state.exit("", nil)
	select {
	case <-cb:
	case <-time.After(time.Duration(1)):
		t.Error("failed to raise onexit")
	}
}
