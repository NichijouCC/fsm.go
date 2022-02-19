package hfsm

type TransitionOptions struct {
	From      string
	To        string
	Condition func(ctx interface{}) bool
}

type Transition struct {
	From      IState
	To        IState
	Condition func(context interface{}) bool
}

func newTransition(from IState, to IState, condition func(ctx interface{}) bool) *Transition {
	return &Transition{
		From:      from,
		To:        to,
		Condition: condition,
	}
}
