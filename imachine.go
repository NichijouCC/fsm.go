package hfsm

import "time"

type IMachine interface {
	GetContext() interface{}
	HasState(state string) bool
	AddState(state IState)
	AddStates(states []IState)
	RemoveState(name string)
	GetCurrent() string
	ChangToState(name string)
	AddTransition(from, to string, condition func(context interface{}) bool)
	RemoveTransition(from, to string)
}

type IState interface {
	GetName() string
	setMachine(machine IMachine)
	addTransition(to IState, condition func(context interface{}) bool)
	removeTransitionTo(to IState)
	hasTransition(to IState) bool
	enter(pre string, ctx interface{})
	update(deltaTime time.Duration, ctx interface{})
	exit(next string, ctx interface{})
}
