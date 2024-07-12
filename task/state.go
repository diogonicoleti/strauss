package task

// TODO: Hard-coded state transitions, we need to use a state machine
// to deal with state transitions
var stateTransitionMap = map[State][]State{
	Pending:   {Scheduled},
	Scheduled: {Scheduled, Running, Failed},
	Running:   {Running, Completed, Failed},
	Completed: {},
	Failed:    {},
}

func ValidStateTransition(curr State, next State) bool {
	for _, s := range stateTransitionMap[curr] {
		if s == next {
			return true
		}
	}
	return false
}
