package process

import (
	"testing"
	"time"
)

func TestCreateAction(t *testing.T) {
	in := make(chan *Event)

	action := NewAction(in, func(out *Event) *Event {
		return nil
	})

	if action.function == nil {
		t.Error("Function should be initialized")
	}

	if action.out == nil {
		t.Error("Out channel should be initialized")
	}

	if action.Out() == nil {
		t.Error("Out method should not return nil")
	}

	if action.in != in {
		t.Error("in channel should be the one provided")
	}

	if action.Done() != nil {
		t.Error("Soure nodes sould not return done channels")
	}
}

func TestStartAction(t *testing.T) {
	functionRan := false

	e := &Event{
		Origin: "none",
	}

	f := func(in *Event) *Event {
		functionRan = true
		in.Origin = "test function"

		return in
	}

	in := make(chan *Event)

	action := NewAction(in, f)
	action.Start()

	out := action.Out()

	time.Sleep(200 * time.Millisecond)

	in <- e
	<-out

	if functionRan == false {
		t.Error("Provided function did not run")
	}

	if e.Origin != "test function" {
		t.Error("Event was not modified")
	}
}
