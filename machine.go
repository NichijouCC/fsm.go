package hfsm

import (
	"log"
	"time"
)

type Machine struct {
	*State
	context interface{}
	name string
	enter *EnterState
	any *AnyState
	states map[string]IState
	current IState
	parent *Machine
}


type machineOption struct {
	parent *Machine
	states []IState
}


func WithState(state []IState) func(opt *machineOption) {
	return func(opt *machineOption) {
		opt.states=state
	}
}

func NewMachine(name string,context interface{},opts ...func(opt *machineOption)) *Machine {
	opt:=&machineOption{}
	for _,el:=range opts{
		el(opt)
	}
	enter:=newEnterState()
	machine:= &Machine{
		State:   NewBaseState(name),
		context: context,
		parent: opt.parent,
		states: map[string]IState{},
		enter:enter,
		current: enter,
	}
	for _,el:=range opt.states{
		machine.AddState(el)
	}
	machine.AddState(enter)
	return machine
}

func (m *Machine) GetContext() interface{} {
	return m.context
}

func (m *Machine) AddEnterToState(toName string,condition func(context interface{}) bool) {
	m.enter.AddTransition(toName,condition)
}

func (m *Machine) RemoveEnterToState(toName string) {
	m.enter.RemoveTransition(toName)
}

func (m *Machine) AddState(state IState)  {
	m.states[state.GetName()]=state
	state.SetMachine(m)
}

func (m *Machine) AddStates(states []IState)  {
	for _,el:=range states{
		m.AddState(el)
	}
}

func (m *Machine) RemoveState(name string)  {
	delete(m.states,name)
}

func (m *Machine) GetState(name string) (IState,bool) {
	stat,ok:=m.states[name]
	return stat,ok
}

func (m *Machine) GetCurrent() IState {
	return m.current
}

func (m *Machine) HasState(state IState) bool {
	if stat,ok:=m.states[state.GetName()];ok{
		return stat==state
	}
	return false
}

func (m *Machine) ChangToState(name string) {
	if m.current.GetName()==name {
		return
	}
	nextState,ok:=m.states[name]
	if !ok{
		log.Printf("warn:failed to changeToState,can't find state %v",name)
		return
	}
	m.current.Exit(nextState)
	nextState.Enter(m.current)
	m.current=nextState
}

func (m *Machine) Update(deltaTime time.Duration)  {
	if m.parent!=nil {
		//parent hierarchy update
		m.State.Update(deltaTime)
		if m.parent.current==m {
			//machine hierarchy any update
			m.any.Update(deltaTime)
			m.current.Update(deltaTime)
		}
	}else {
		//machine hierarchy any update
		m.any.Update(deltaTime)
		m.current.Update(deltaTime)
	}
}