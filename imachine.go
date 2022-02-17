package hfsm

import "time"

type IStateMachine interface {
	GetContext() interface{}
	HasState(state IState) bool
	GetState(name string) (IState,bool)
	AddState(state IState)
	AddStates(states []IState)
	RemoveState(name string)
	GetCurrent()IState
	GetEnterState() *EnterState
	GetAnyState() *AnyState
	ChangToState(name string)
}

type IState interface {
	IExtendState
	GetName() string
	SetMachine(machine IStateMachine)
	GetMachine() IStateMachine
	GetContext() interface{}
	AddTransitionTo(toName string,condition func(context interface{})bool)
	RemoveTransitionTo(toName string)
	Enter(pre IState)
	Update(deltaTime time.Duration)
	Exit(next IState)
}

type IExtendState interface {
	OnEnter(prev IState)
	OnUpdate(deltaTime time.Duration)
	OnExit(next IState)
}