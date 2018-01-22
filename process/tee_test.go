package process

import (
	"testing"
)

func TestTee(t *testing.T) {
	in := make(chan *Event)
	tee := NewTee(in)

	outA := tee.Out()
	outB := tee.Out()

	if outA == outB {
		t.Error("Out channels should be different")
	}

	tee.Start()

	e := &Event{
		Origin: "test",
	}

	in <- e

	eA := <-outA
	eB := <-outB

	if eA != e || eB != e {
		t.Error("Output should be the same as input in: %v, a: %v, b %v",
			e, eA, eB)
	}

	close(in)

	_, okA := <-outA
	_, okB := <-outB

	if okA || okB {
		t.Error("Outputs should be closed when input is also closed a: %v, b: %v",
			okA, okB)
	}

	if tee.Done() != nil {
		t.Error("Done should return nil")
	}
}
