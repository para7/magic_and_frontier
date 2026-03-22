package idseq

type CounterState struct {
	CastIDs int `json:"castids"`
}

func EmptyCounterState() CounterState {
	return CounterState{}
}

func NextCastID(state CounterState) (CounterState, int) {
	state.CastIDs++
	return state, state.CastIDs
}
