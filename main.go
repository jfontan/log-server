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

	tee := process.NewTee(filter.Out())
	tee.Start()

	printRootSink := process.GenPrintKeySink("root")
	printMsgSink := process.GenPrintKeySink("msg")

	sink1 := process.NewSink(tee.Out(), printMsgSink)
	sink1.Start()

	sink2 := process.NewSink(tee.Out(), printRootSink)
	sink2.Start()

	<-sink1.Done()
	<-sink2.Done()
}
