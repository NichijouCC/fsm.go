package hfsm

import (
	"log"
	"time"
)

type State struct {
	name string
	machine IStateMachine
	transitions []*Transition
	enterTime time.Time
	duration time.Duration
}

func NewBaseState(name string) *State {
	return &State{
		name: name,
		transitions: []*Transition{},
	}
}

func (b *State) SetMachine(machine IStateMachine) {
	b.machine=machine
}

func (b *State) GetMachine() IStateMachine {
	return b.machine
}

func (b *State) GetContext() interface{} {
	if b.machine==nil {
		return nil
	}
	return b.machine.GetContext()
}

func (b *State) GetName() string {
	return b.name
}

func (b *State) AddTransition(toName string,condition func(context interface{}) bool) {
	if b.machine==nil {
		return
	}
	if to,ok:=b.machine.GetState(toName);ok {
		for _,el:=range b.transitions{
			if el.To.GetName()==toName {
				log.Printf("warn:failed to add Transition, already contain it(%v->%v)",b.GetName(),toName)
				return
			}
		}
		b.transitions=append(b.transitions,newTransition(b,to,condition))
	}else {
		log.Printf("warn:failed to add Transition(%v->%v), can't find state %v",b.GetName(),toName,toName)
	}
}

func (b *State) RemoveTransition(toName string) {
	for index,el:=range b.transitions{
		if el.To.GetName()==toName {
			b.transitions=append(b.transitions[:index], b.transitions[index+1:]...)
			return
		}
	}
	log.Printf("warn:failed to reomve transition, can't find it(%v->%v)",b.GetName(),toName)
}

func (b *State) Update(deltaTime time.Duration) {
	for _,el:=range b.transitions{
		if el.Condition==nil||(el.Condition!=nil&&el.Condition(b.GetContext())) {
			b.machine.ChangToState(el.To.GetName())
			return
		}
	}
	b.duration=time.Now().Sub(b.enterTime)
	b.OnUpdate(deltaTime)
}

func (b *State) Enter(pre IState) {
	b.enterTime=time.Now()
	b.OnEnter(pre)
}

func (b *State) Exit(next IState) {
	b.OnExit(next)
}

func (b *State) OnEnter(prev IState) {

}

func (b *State) OnExit(next IState) {

}

func (b *State) OnUpdate(deltaTime time.Duration) {

}