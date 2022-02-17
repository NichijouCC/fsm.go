package hfsm

import (
	"log"
	"time"
)

type BaseState struct {
	Extend IExtendState
	name string
	machine IStateMachine
	transitions []*Transition
	enterTime time.Time
	duration time.Duration
}

func NewBaseState(name string) *BaseState {
	return &BaseState{
		name: name,
		transitions: []*Transition{},
	}
}

func (b *BaseState) SetMachine(machine IStateMachine) {
	b.machine=machine
}

func (b *BaseState) GetMachine() IStateMachine {
	return b.machine
}

func (b *BaseState) GetContext() interface{} {
	if b.machine==nil {
		return nil
	}
	return b.machine.GetContext()
}

func (b *BaseState) GetName() string {
	return b.name
}

func (b *BaseState) AddTransitionTo(toName string,condition func(context interface{}) bool) {
	if b.machine==nil {
		return
	}
	if to,ok:=b.machine.GetState(toName);ok {
		for _,el:=range b.transitions{
			if el.To==to {
				log.Printf("warn:failed to add Transition, already contain it(%v->%v)",b.GetName(),toName)
				return
			}
		}
		b.transitions=append(b.transitions,newTransition(b,to,condition))
	}else {
		log.Printf("warn:failed to add Transition(%v->%v), can't find state %v",b.GetName(),toName,toName)
	}
}

func (b *BaseState) RemoveTransitionTo(toName string) {
	for index,el:=range b.transitions{
		if el.To.GetName()==toName {
			b.transitions=append(b.transitions[:index], b.transitions[index+1:]...)
			return
		}
	}
	log.Printf("warn:failed to reomve transition, can't find it(%v->%v)",b.GetName(),toName)
}

func (b *BaseState) Update(deltaTime time.Duration) {
	for _,el:=range b.transitions{
		if el.Condition==nil||(el.Condition!=nil&&el.Condition(b.GetContext())) {
			b.machine.ChangToState(el.To.GetName())
			return
		}
	}
	b.duration=time.Now().Sub(b.enterTime)
	b.Extend.OnUpdate(deltaTime)
}

func (b *BaseState) Enter(pre IState) {
	b.enterTime=time.Now()
	b.Extend.OnEnter(pre)
}

func (b *BaseState) Exit(next IState) {
	b.Extend.OnExit(next)
}

func (b *BaseState) OnEnter(prev IState) {

}

func (b *BaseState) OnExit(next IState) {

}

func (b *BaseState) OnUpdate(deltaTime time.Duration) {

}