package hfsm

type Transition struct {
	From IState
	To IState
	Condition func(context interface{})bool
}

func newTransition(from IState, to IState, condition func(context interface{})bool) *Transition {
	return &Transition{
		From: from,
		To: to,
		Condition: condition,
	}
}