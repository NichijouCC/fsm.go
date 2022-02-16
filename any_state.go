package hfsm

type AnyState struct {
	*State
}

func NewAnyState() *AnyState {
	return &AnyState{
		NewBaseState("any"),
	}
}