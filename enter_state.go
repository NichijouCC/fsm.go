package hfsm

type EnterState struct {
	*BaseState
}

func newEnterState() *EnterState {
	return &EnterState{
		NewBaseState(ENTER),
	}
}