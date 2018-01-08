package process

import "regexp"

type BasicNode interface {
	Start()
	Out() chan *Event
	Done() chan bool
}

type ProcessFunc func(in *Event) (out *Event)

type Process struct {
	function ProcessFunc
	in       chan *Event
	out      chan *Event
}

func NewProcess(in chan *Event, f ProcessFunc) *Process {
	out := make(chan *Event)

	return &Process{
		function: f,
		in:       in,
		out:      out,
	}
}

func (s *Process) Start() {
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

func (s *Process) Out() chan *Event {
	return s.out
}

func (s *Process) Done() chan bool {
	return nil
}

func GenRegexpFilter(filter map[string]string) ProcessFunc {
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
