package process

type Tee struct {
	in   chan *Event
	outs []chan *Event
}

func NewTee(in chan *Event) *Tee {
	return &Tee{
		in:   in,
		outs: make([]chan *Event, 0),
	}
}

func (t *Tee) Start() {
	go func() {
		for event := range t.in {
			for _, out := range t.outs {
				out <- event
			}
		}

		for _, out := range t.outs {
			close(out)
		}
	}()
}

func (t *Tee) Out() chan *Event {
	out := make(chan *Event)
	t.outs = append(t.outs, out)

	return out
}
