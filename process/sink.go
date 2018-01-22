package process

import "fmt"

type SinkFunc func(in *Event)

type Sink struct {
	function SinkFunc
	in       chan *Event
	done     chan bool
}

func NewSink(in chan *Event, f SinkFunc) *Sink {
	done := make(chan bool)

	return &Sink{
		function: f,
		in:       in,
		done:     done,
	}
}

func (s *Sink) Start() {
	go func() {
		for event := range s.in {
			s.function(event)
		}

		s.done <- true
	}()
}

func (s *Sink) Done() chan bool {
	return s.done
}

func (s *Sink) Out() chan *Event {
	return nil
}

func GenPrintKeySink(key string) SinkFunc {
	return func(event *Event) {
		fmt.Println(event.Data[key])
	}
}

func GenPrintCounters() SinkFunc {
	return func(event *Event) {
		for name, count := range counters {
			fmt.Printf("%v = %v\n", name, count)
		}
	}
}
