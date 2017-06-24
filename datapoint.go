package main

type Datapoint struct {
	name  string
	score float64
}

type DatapointSort []*Datapoint

func (dps DatapointSort) Len() int {
	return len(dps)
}

func (dps DatapointSort) Swap(i, j int) {
	dps[i], dps[j] = dps[j], dps[i]
}

func (dps DatapointSort) Less(i, j int) bool {
	return dps[i].score < dps[j].score
}
