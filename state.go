package hfsm

import (
	"log"
	"time"
)

type State struct {
	name        string
	machine     IMachine
	transitions []*Transition

	OnEnter  func(pre string, ctx interface{})
	OnUpdate func(deltaTime time.Duration, ctx interface{})
	OnExit   func(next string, ctx interface{})
}

func (b *State) setMachine(machine IMachine) {
	b.machine = machine
}

type StateOptions struct {
	OnEnter  func(pre string, ctx interface{})
	OnUpdate func(deltaTime time.Duration, ctx interface{})
	OnExit   func(next string, ctx interface{})
}

func WithOnEnter(enter func(pre string, ctx interface{})) func(opt *StateOptions) {
	return func(opt *StateOptions) {
		opt.OnEnter = enter
	}
}

func WithOnUpdate(update func(deltaTime time.Duration, ctx interface{})) func(opt *StateOptions) {
	return func(opt *StateOptions) {
		opt.OnUpdate = update
	}
}

func WithOnExit(exit func(next string, ctx interface{})) func(opt *StateOptions) {
	return func(opt *StateOptions) {
		opt.OnExit = exit
	}
}

func NewState(name string, opts ...func(opt *StateOptions)) *State {
	opt := &StateOptions{}
	for _, el := range opts {
		el(opt)
	}
	state := &State{
		name:        name,
		transitions: []*Transition{},
	}
	state.OnEnter = opt.OnEnter
	state.OnUpdate = opt.OnUpdate
	state.OnExit = opt.OnExit
	return state
}

func (b *State) GetName() string {
	return b.name
}

func (b *State) addTransition(to IState, condition func(context interface{}) bool) {
	for _, el := range b.transitions {
		if el.To == to {
			log.Printf("warn:failed to add Transition, already contain it(%v->%v)", b.name, to.GetName())
			return
		}
	}
	b.transitions = append(b.transitions, newTransition(b, to, condition))
}

func (b *State) removeTransitionTo(to IState) {
	for index, el := range b.transitions {
		if el.To == to {
			b.transitions = append(b.transitions[:index], b.transitions[index+1:]...)
			return
		}
	}
	log.Printf("warn:failed to reomve transition, can't find it(%v->%v)", b.name, to.GetName())
}

func (b *State) hasTransition(to IState) bool {
	for _, el := range b.transitions {
		if el.To == to {
			return true
		}
	}
	return false
}

func (b *State) update(deltaTime time.Duration, ctx interface{}) {
	for _, el := range b.transitions {
		if el.Condition == nil || (el.Condition != nil && el.Condition(ctx)) {
			b.machine.ChangToState(el.To.GetName())
			return
		}
	}
	if b.OnUpdate != nil {
		b.OnUpdate(deltaTime, ctx)
	}
}

func (b *State) enter(pre string, ctx interface{}) {
	if b.OnEnter != nil {
		b.OnEnter(pre, ctx)
	}
}

func (b *State) exit(next string, ctx interface{}) {
	if b.OnExit != nil {
		b.OnExit(next, ctx)
	}
}
