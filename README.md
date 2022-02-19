# HFSM
Hierarchical Finite State Machine Implementation in Go

Simple Example
```
ctx:= struct {}{}
machine:=NewMachine(
    "test",
    ctx,
    WithStates(map[string]*StateOptions{
        "state_a":&StateOptions{OnEnter: func(pre string, ctx interface{}) {
            log.Println("state a enter")
    }}}),
    WithTransitions([]*TransitionOptions{&TransitionOptions{From: ENTER,To:"state_a"}}))
stateB:=NewState("state_b",WithOnEnter(func(pre string, ctx interface{}) {
    log.Println("state b enter")
}))
stateB.OnExit= func(next string, ctx interface{}) {
    log.Println("state b exit")
}
machine.AddState(stateB)
subMachine:=NewMachine("subMachine",ctx)
subMachine.OnEnter= func(pre string, ctx interface{}) {
    log.Println("subMachine enter")
}
machine.AddState(subMachine)
machine.AddTransition("state_a","subMachine",nil)

machine.Update(time.Duration(1))

```
