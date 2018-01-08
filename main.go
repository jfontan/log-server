package main

import (
	"github.com/jfontan/log-server/process"
)

func main() {
	f := map[string]string{
		"lvl":   "^eror$",
		"error": "index read failed",
	}

	fileSource := process.GenFileSource("logs")
	source := process.NewSource(fileSource)
	source.Start()

	regexpFilter := process.GenRegexpFilter(f)
	filter := process.NewProcess(source.Out(), regexpFilter)
	filter.Start()

	printSink := process.GenPrintKeySink("root")
	sink := process.NewSink(filter.Out(), printSink)
	sink.Start()

	<-sink.Done()
}
