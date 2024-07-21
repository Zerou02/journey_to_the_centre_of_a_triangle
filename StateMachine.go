package main

type State interface {
	Process(delta float32)
	Draw()
	IsFinished() bool
}

type StateMachine struct {
	states    []State
	currState int
}

func newStateMachine() StateMachine {
	return StateMachine{currState: 0}
}

func (this *StateMachine) addState(state State) *StateMachine {
	this.states = append(this.states, state)
	return this
}

func (this *StateMachine) process(delta float32) {
	if this.isFinished() {
		return
	}
	var currState = this.states[this.currState]
	currState.Process(delta)
	if currState.IsFinished() {
		this.currState++
	}
}

func (this *StateMachine) draw() {
	if this.isFinished() {
		return
	}
	var currState = this.states[this.currState]
	currState.Draw()
}

func (this *StateMachine) isFinished() bool {
	return this.currState >= len(this.states)
}
