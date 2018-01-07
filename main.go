package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/jfontan/log-server/parse"
)

type StringMap map[string]string

type Event struct {
	Origin string
	Data   StringMap
}

type GeneratorFunc func(source string) (out chan *Event)

type ProcessFunc func(in chan *Event) (out chan *Event)

type SinkFun func(in chan *Event) (done chan bool)

func GenerateFromFile(fileName string) chan *Event {
	out := make(chan *Event)
	repeat := false

	go func() {
		for firstTime := true; firstTime || repeat; firstTime = false {
			file, err := os.Open(fileName)
			if err != nil {
				panic(err)
			}

			scanner := bufio.NewScanner(file)

			// Increase max token size as some lines are longer that 64k
			maxTokenSize := 16 * 1024 * 1024
			scanner.Buffer(make([]byte, maxTokenSize), maxTokenSize)

			for scanner.Scan() {
				parsed := parse.ParseLine(scanner.Text())
				if parsed != nil {
					out <- &Event{Origin: fileName, Data: parsed}
				}
			}

			print(scanner.Text())

			file.Close()
		}
		close(out)
	}()

	return out
}

func SinkStdout(in chan *Event) chan bool {
	done := make(chan bool)

	go func() {
		for event := range in {
			fmt.Println(event.Data)
		}

		done <- true
	}()

	return done
}

func SinkStdoutKey(key string, in chan *Event) chan bool {
	done := make(chan bool)

	go func() {
		for event := range in {
			fmt.Println(event.Data[key])
		}

		done <- true
	}()

	return done
}

func FilterRegexp(filter map[string]string, in chan *Event) chan *Event {
	out := make(chan *Event)

	regs := make(map[string]*regexp.Regexp)
	for k, v := range filter {
		regs[k] = regexp.MustCompile(v)
	}

	total := 0
	filtered := 0

	go func() {
		for event := range in {
			valid := true
			total += 1

			for k, v := range regs {
				val, ok := event.Data[k]

				if !ok {
					valid = false
					break
				}

				if !v.MatchString(val) {
					valid = false
					break
				}
			}

			if valid {
				filtered += 1
				out <- event
			}
		}

		fmt.Printf("Total: %v, Filtered: %v\n", total, filtered)

		close(out)
	}()

	return out
}

func printError(c chan StringMap) {
	for msg := range c {
		if msg["lvl"] == "eror" {
			fmt.Printf("%q\n", msg)
		}
	}
}

func main() {
	f := map[string]string{
		"lvl":   "^eror$",
		"error": "index read failed",
	}

	generator := GenerateFromFile("logs")
	filter := FilterRegexp(f, generator)
	done := SinkStdoutKey("root", filter)

	<-done

	// time.Sleep(1000 * time.Second)
}
