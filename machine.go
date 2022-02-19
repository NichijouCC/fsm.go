package hfsm

import (
	"log"
	"time"
)

const (
	ENTER = "HFSM_ENTER"
	ANY   = "HFSM_ANY"
)

type Machine struct {
	IState
	context    interface{}
	enterState IState
	anyState   IState
	states     map[string]IState
	current    IState
	parent     *Machine

	OnEnter  func(pre string, ctx interface{})
	OnUpdate func(deltaTime time.Duration, ctx interface{})
	OnExit   func(next string, ctx interface{})
}

type machineOption struct {
	States          map[string]*StateOptions
	Transitions     []*TransitionOptions
	OnMachineEnter  func(pre string, ctx interface{})
	OnMachineUpdate func(deltaTime time.Duration, ctx interface{})
	OnMachineExit   func(next string, ctx interface{})
}

func WithStates(states map[string]*StateOptions) func(opt *machineOption) {
	return func(opt *machineOption) {
		opt.States = states
	}
}

func WithTransitions(transitions []*TransitionOptions) func(opt *machineOption) {
	return func(opt *machineOption) {
		opt.Transitions = transitions
	}
}

func WithMachineOnEnter(enter func(pre string, ctx interface{})) func(opt *machineOption) {
	return func(opt *machineOption) {
		opt.OnMachineEnter = enter
	}
}

func WithMachineOnUpdate(update func(deltaTime time.Duration, ctx interface{})) func(opt *machineOption) {
	return func(opt *machineOption) {
		opt.OnMachineUpdate = update
	}
}

func WithMachineOnExit(exit func(next string, ctx interface{})) func(opt *machineOption) {
	return func(opt *machineOption) {
		opt.OnMachineExit = exit
	}
}

func NewMachine(name string, context interface{}, opts ...func(opt *machineOption)) *Machine {
	opt := &machineOption{}
	for _, el := range opts {
		el(opt)
	}
	enter := NewState(ENTER)
	any := NewState(ANY)
	machine := &Machine{
		context:    context,
		states:     map[string]IState{},
		enterState: enter,
		anyState:   any,
		current:    enter,
	}
	state := NewState(name)
	machine.IState = state
	machine.OnEnter = opt.OnMachineEnter
	machine.OnUpdate = opt.OnMachineUpdate
	machine.OnExit = opt.OnMachineExit

	machine.AddState(enter)
	machine.AddState(any)
	for stateName, el := range opt.States {
		machine.AddState(NewState(stateName, WithOnEnter(el.OnEnter), WithOnUpdate(el.OnUpdate), WithOnExit(el.OnExit)))
	}
	for _, el := range opt.Transitions {
		machine.AddTransition(el.From, el.To, el.Condition)
	}
	return machine
}

func (m *Machine) GetContext() interface{} {
	return m.context
}

func (m *Machine) AddState(state IState) {
	m.states[state.GetName()] = state
	state.setMachine(m)
}

func (m *Machine) AddStates(states []IState) {
	for _, el := range states {
		m.AddState(el)
	}
}

func (m *Machine) RemoveState(name string) {
	delete(m.states, name)
}

func (m *Machine) GetCurrent() string {
	return m.current.GetName()
}

func (m *Machine) HasState(state string) bool {
	return m.states[state] != nil
}

func (m *Machine) AddTransition(from, to string, condition func(context interface{}) bool) {
	fromState, ok := m.states[from]
	if !ok {
		log.Printf("warn:failed to add Transition(%v->%v), can't find state %v", from, to, from)
	}
	toState, ok := m.states[to]
	if !ok {
		log.Printf("warn:failed to add Transition(%v->%v), can't find state %v", from, to, to)
	}
	fromState.addTransition(toState, condition)
}

func (m *Machine) RemoveTransition(from, to string) {
	fromState, ok := m.states[from]
	if !ok {
		log.Printf("warn:failed to remove Transition(%v->%v), can't find state %v", from, to, from)
	}
	toState, ok := m.states[to]
	if !ok {
		log.Printf("warn:failed to remove Transition(%v->%v), can't find state %v", from, to, to)
	}
	fromState.removeTransitionTo(toState)
}

func (m *Machine) HasTransition(from, to string) bool {
	fromState, ok := m.states[from]
	if !ok {
		return false
	}
	toState, ok := m.states[to]
	if !ok {
		return false
	}
	return fromState.hasTransition(toState)
}

func (m *Machine) enter(pre string, ctx interface{}) {
	m.current = m.enterState
	if m.OnEnter != nil {
		m.OnEnter(pre, ctx)
	}
}

func (m *Machine) exit(next string, ctx interface{}) {
	if m.OnExit != nil {
		m.OnExit(next, ctx)
	}
}

func (m *Machine) update(deltaTime time.Duration, ctx interface{}) {
	m.IState.update(deltaTime, ctx)
	m.Update(deltaTime)
	if m.OnUpdate != nil {
		m.OnUpdate(deltaTime, ctx)
	}
}

func (m *Machine) ChangToState(name string) {
	if m.current.GetName() == name {
		return
	}
	nextState, ok := m.states[name]
	if !ok {
		log.Printf("warn:failed to changeToState,can't find state %v", name)
		return
	}
	m.current.exit(name, m.context)
	nextState.enter(m.current.GetName(), m.context)
	m.current = nextState
}

//machine internal update
func (m *Machine) Update(deltaTime time.Duration) {
	current := m.current.GetName()
	m.anyState.update(deltaTime, m.context)
	m.current.update(deltaTime, m.context)
	if m.current.GetName() != current {
		m.Update(deltaTime)
	}
}
