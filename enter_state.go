package hfsm

type EnterState struct {
	*State
}

func newEnterState() *EnterState {
	return &EnterState{
		NewBaseState("enter"),
	}
}