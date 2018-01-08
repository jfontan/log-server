package process

import (
	"bufio"
	"os"

	"github.com/jfontan/log-server/parse"
)

type SourceFunc func(out chan *Event)

type Source struct {
	function SourceFunc
	out      chan *Event
}

func NewSource(f SourceFunc) (source *Source) {
	out := make(chan *Event)

	return &Source{
		function: f,
		out:      out,
	}
}

func (s *Source) Start() {
	go s.function(s.out)
}

func (s *Source) Out() chan *Event {
	return s.out
}

func GenFileSource(fileName string) SourceFunc {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	// Increase max token size as some lines are longer that 64k
	maxTokenSize := 16 * 1024 * 1024
	scanner.Buffer(make([]byte, maxTokenSize), maxTokenSize)

	return func(out chan *Event) {
		defer file.Close()
		defer close(out)

		for scanner.Scan() {
			parsed := parse.ParseLine(scanner.Text())
			if parsed != nil {
				out <- &Event{Origin: fileName, Data: parsed}
			}
		}
	}
}
