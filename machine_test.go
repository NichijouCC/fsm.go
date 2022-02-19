package hfsm

import (
	"testing"
	"time"
)

func TestNewMachine(t *testing.T) {
	ctx := struct{}{}
	machine := NewMachine("test", ctx)
	if machine.GetName() != "test" {
		t.Error("machine name not right")
	}
	if machine.GetCurrent() != ENTER {
		t.Error("machine init state not be ", ENTER)
	}
	if machine.HasState(ANY) == false {
		t.Error("machine not contain anyState")
	}
	if machine.HasState(ENTER) == false {
		t.Error("machine not contain enterState")
	}
}

func TestWithStates(t *testing.T) {
	ctx := struct{}{}
	states := map[string]*StateOptions{"idle": &StateOptions{}}
	machine := NewMachine("test", ctx, WithStates(states))
	if machine.HasState("idle") == false {
		t.Error("failed to add idle state")
	}
}

func TestWithTransitions(t *testing.T) {
	ctx := struct{}{}
	states := map[string]*StateOptions{"idle": &StateOptions{}}
	transitions := []*TransitionOptions{&TransitionOptions{From: ENTER, To: "idle"}}
	machine := NewMachine("test", ctx, WithStates(states), WithTransitions(transitions))
	if machine.HasTransition(ENTER, "idle") == false {
		t.Error("failed to add transition")
	}
}

func TestWithMachineOnEnter(t *testing.T) {
	ctx := struct{}{}
	cbChan := make(chan struct{}, 1)
	machine := NewMachine("test", ctx, WithMachineOnEnter(func(pre string, ctx interface{}) {
		cbChan <- struct{}{}
	}))
	machine.enter("", ctx)
	select {
	case <-cbChan:
	case <-time.After(time.Second):
		t.Error("failed to raise 'onEnter'")
	}
}

func TestWithMachineOnUpdate(t *testing.T) {
	ctx := struct{}{}
	cbChan := make(chan struct{}, 1)
	machine := NewMachine("test", ctx, WithMachineOnUpdate(func(deltaTime time.Duration, ctx interface{}) {
		cbChan <- struct{}{}
	}))
	machine.update(time.Duration(1), ctx)
	select {
	case <-cbChan:
	case <-time.After(time.Second):
		t.Error("failed to raise 'onUpdate'")
	}
}

func TestWithMachineOnExit(t *testing.T) {
	ctx := struct{}{}
	cbChan := make(chan struct{}, 1)
	machine := NewMachine("test", ctx, WithMachineOnExit(func(pre string, ctx interface{}) {
		cbChan <- struct{}{}
	}))
	machine.exit("", ctx)
	select {
	case <-cbChan:
	case <-time.After(time.Second):
		t.Error("failed to raise 'onExit'")
	}
}

func TestMachine_AddState(t *testing.T) {
	machine := NewMachine("test", nil)
	state := NewState("idle")
	machine.AddState(state)
	if machine.HasState("idle") == false {
		t.Error("failed to add state idle")
	}
}

func TestMachine_RemoveState(t *testing.T) {
	machine := NewMachine("test", nil)
	state := NewState("idle")
	machine.AddState(state)
	machine.RemoveState("idle")
	if machine.HasState("idle") {
		t.Error("failed to remove state idle")
	}
}

func TestMachine_HasState(t *testing.T) {
	machine := NewMachine("test", nil)
	state := NewState("idle")
	machine.AddState(state)
	machine.AddTransition(ENTER, "idle", nil)
	if machine.HasTransition(ENTER, "idle") == false {
		t.Error("failed to add transition")
	}
}

func TestMachine_AddTransition(t *testing.T) {
	machine := NewMachine("test", nil)
	state := NewState("idle")
	machine.AddState(state)
	machine.AddTransition(ENTER, "idle", nil)
	if machine.HasTransition(ENTER, "idle") == false {
		t.Error("failed to add transition")
	}
}

func TestMachine_HasTransition(t *testing.T) {
	machine := NewMachine("test", nil)
	if machine.HasTransition(ENTER, "idle") {
		t.Error("failed to check whether have transition")
	}
}

func TestMachine_RemoveTransition(t *testing.T) {
	machine := NewMachine("test", nil)
	state := NewState("idle")
	machine.AddState(state)
	machine.AddTransition(ENTER, "idle", nil)
	machine.RemoveTransition(ENTER, "idle")
	if machine.HasTransition(ENTER, "idle") {
		t.Error("failed to remove transition")
	}
}

func TestMachine_ChangToState(t *testing.T) {
	machine := NewMachine("test", nil)
	state := NewState("idle")
	cb := make(chan struct{}, 1)
	state.OnEnter = func(pre string, ctx interface{}) {
		cb <- struct{}{}
	}
	machine.AddState(state)
	machine.ChangToState("idle")
	if machine.GetCurrent() != "idle" {
		t.Error("failed to change state to idle")
	}
	select {
	case <-cb:
	case <-time.After(time.Duration(1)):
		t.Error("failed to raise enter state OnEnter when machine changeState")
	}
	exitCb := make(chan struct{}, 1)
	state.OnExit = func(next string, ctx interface{}) {
		exitCb <- struct{}{}
	}
	machine.ChangToState(ENTER)
	if machine.GetCurrent() != ENTER {
		t.Error("failed to raise exit state OnExit when machine changeState")
	}
}

func TestMachine_Update(t *testing.T) {
	machine := NewMachine("test", nil)
	state := NewState("idle")
	machine.AddState(state)
	machine.AddTransition(ENTER, "idle", nil)
	machine.Update(time.Duration(1))
	if machine.GetCurrent() != "idle" {
		t.Error("failed to chang state by 'Update' when condition is fulfilled")
	}
}

func TestMachine_subMachine_enterAndOutEvent(t *testing.T) {
	machine := NewMachine("test", nil)
	subMachine := NewMachine("sub", nil)
	cb1 := make(chan struct{}, 1)
	subMachine.OnEnter = func(pre string, ctx interface{}) {
		if pre == ENTER {
			cb1 <- struct{}{}
		}
	}
	machine.AddState(subMachine)
	machine.ChangToState("sub")
	select {
	case <-cb1:
	case <-time.After(time.Duration(1)):
		t.Error("subMachine onEnter not raise")
	}

	cb2 := make(chan struct{}, 1)
	subMachine.OnExit = func(next string, ctx interface{}) {
		if next == ENTER {
			cb2 <- struct{}{}
		}
	}
	machine.ChangToState(ENTER)
	select {
	case <-cb2:
	case <-time.After(time.Duration(1)):
		t.Error("subMachine onExit not raise")
	}
}

func TestMachine_subMachine_update(t *testing.T) {
	machine := NewMachine("test", nil)
	subMachine := NewMachine("sub", nil)
	cb := make(chan struct{}, 1)
	subMachine.OnUpdate = func(deltaTime time.Duration, ctx interface{}) {
		cb <- struct{}{}
	}
	machine.AddState(subMachine)
	machine.ChangToState("sub")

	state := NewState("idle")
	subMachine.AddState(state)
	subMachine.AddTransition(ENTER, "idle", nil)
	machine.Update(time.Duration(1))
	select {
	case <-cb:
	case <-time.After(time.Duration(1)):
		t.Error("subMachine onUpdate not raise")
	}
	if subMachine.GetCurrent() != "idle" {
		t.Error("failed to chang subMachine state by 'Update' when condition is fulfilled")
	}
}

func TestMachine_subMachine_ChangeState(t *testing.T) {
	machine := NewMachine("test", nil)
	subMachine := NewMachine("sub", nil)
	machine.AddState(subMachine)
	state := NewState("idle")
	subMachine.AddState(state)
	subMachine.AddTransition(ENTER, "idle", nil)
	machine.ChangToState("sub")
	machine.AddTransition("sub", ENTER, nil)
	machine.Update(time.Duration(1))
	if machine.GetCurrent() != ENTER {
		t.Error("failed to change Machine state by 'Update' when condition is fulfilled")
	}
}
