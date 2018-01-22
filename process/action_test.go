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

func TestRegexpFilter(t *testing.T) {
	filter := map[string]string{
		"one": "one",
		"two": "^two",
	}

	cases := []struct {
		Event *Event
		Nil   bool
	}{
		{
			Event: &Event{
				Data: map[string]string{
					"one": "one",
				},
			},
			Nil: true,
		},
		{
			Event: &Event{
				Data: map[string]string{
					"one": "one",
					"two": "two",
				},
			},
			Nil: false,
		},
		{
			Event: &Event{
				Data: map[string]string{
					"two": "two",
				},
			},
			Nil: true,
		},
	}

	f := GenRegexpFilter(filter)

	for _, test := range cases {
		res := f(test.Event)

		if (res == nil) != test.Nil {
			if test.Nil {
				t.Errorf("The event %v should not past the test", test.Event)
			} else {
				t.Errorf("The event %v should past the test", test.Event)
			}
		}
	}
}

func testCounter(t *testing.T, name string, expected uint) {
	val, ok := counters[name]

	if !ok {
		t.Errorf("Counter %v should be defined", name)
	}

	if val != expected {
		t.Errorf("Counter %v has value %v, is expected to be %v",
			name, val, expected)
	}
}

func TestCounterAction(t *testing.T) {
	if counters != nil {
		t.Fatal("Counter must be nil when the test starts")
	}

	counterA := GenCounterAction("a")
	counterB := GenCounterAction("b")

	event := &Event{
		Origin: "testing",
	}

	if counterA(event) != event {
		t.Error("Counter function should return the same event")
	}

	counterA(nil)
	counterB(nil)

	testCounter(t, "a", 2)
	testCounter(t, "b", 1)
}
