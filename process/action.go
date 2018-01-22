package process

import (
	"regexp"
)

type ActionFunc func(in *Event) (out *Event)

type Action struct {
	function ActionFunc
	in       chan *Event
	out      chan *Event
}

func NewAction(in chan *Event, f ActionFunc) *Action {
	out := make(chan *Event)

	return &Action{
		function: f,
		in:       in,
		out:      out,
	}
}

func (s *Action) Start() {
	go func() {
		defer close(s.out)

		for event := range s.in {
			out := s.function(event)

			if out != nil {
				s.out <- out
			}
		}
	}()
}

func (s *Action) Out() chan *Event {
	return s.out
}

func (s *Action) Done() chan bool {
	return nil
}

func GenRegexpFilter(filter map[string]string) ActionFunc {
	regs := make(map[string]*regexp.Regexp)
	for k, v := range filter {
		regs[k] = regexp.MustCompile(v)
	}

	return func(event *Event) *Event {
		for k, v := range regs {
			val, ok := event.Data[k]

			if !ok {
				return nil
			}

			if !v.MatchString(val) {
				return nil
			}
		}

		return event
	}
}

var counters map[string]uint

func GenCounterAction(name string) ActionFunc {
	if counters == nil {
		counters = make(map[string]uint)
	}

	_, ok := counters[name]
	if ok {
		panic("Trying to redefine counter named " + name)
	}

	counters[name] = 0

	return func(event *Event) *Event {
		counters[name]++
		return event
	}
}
