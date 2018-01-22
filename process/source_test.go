package process

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestCreateSource(t *testing.T) {
	source := NewSource(func(out chan *Event) {})

	if source.function == nil {
		t.Error("Out channel should be initialized")
	}

	if source.out == nil {
		t.Error("Out channel should be initialized")
	}

	if source.Out() == nil {
		t.Error("Out method should not return nil")
	}

	if source.Done() != nil {
		t.Error("Soure nodes sould not return done channels")
	}
}

func TestStartSource(t *testing.T) {
	functionRan := false

	f := func(out chan *Event) {
		functionRan = true
	}

	source := NewSource(f)
	source.Start()

	time.Sleep(200 * time.Millisecond)

	if functionRan == false {
		t.Error("Provided function did not run")
	}
}

var (
	exampleLogFile = `a=a b=1 c="hello world"`
	exampleLogData = map[string]string{
		"a": "a",
		"b": "1",
		"c": "hello world",
	}
)

func TestGenFileSource(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.WriteString(exampleLogFile); err != nil {
		t.Fatal(err)
	}

	f := GenFileSource(tmpFile.Name())

	out := make(chan *Event)

	go f(out)

	e := <-out

	for k, v := range exampleLogData {
		if e.Data[k] != v {
			t.Errorf("Parsed key %v has incorrect value %v, expected %v",
				k, e.Data[k], v)
		}
	}
}
