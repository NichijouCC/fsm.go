package hfsm

type AnyState struct {
	*BaseState
}

func NewAnyState() *AnyState {
	return &AnyState{
		NewBaseState(ANY),
	}
}