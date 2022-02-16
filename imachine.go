package hfsm

import "time"

type IStateMachine interface {
	GetContext() interface{}
	AddEnterToState(state string,condition func(context interface{})bool)
	RemoveEnterToState(state string)
	HasState(state IState) bool
	GetState(name string) (IState,bool)
	ChangToState(state string)
}

type IState interface {
	GetName() string
	SetMachine(machine IStateMachine)
	GetMachine() IStateMachine
	GetContext() interface{}
	AddTransition(toName string,condition func(context interface{})bool)
	RemoveTransition(toName string)
	OnEnter(prev IState)
	OnUpdate(deltaTime time.Duration)
	OnExit(next IState)
	Enter(pre IState)
	Update(deltaTime time.Duration)
	Exit(next IState)
}
