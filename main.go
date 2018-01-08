package main

import (
	"github.com/jfontan/log-server/process"
)

func main() {
	nodes := process.LoadNodes("process.yml")
	done := nodes.Start()

	// wait for all sinks to finish
	for _, d := range done {
		<-d
	}
}
