package main

import (
	"github.com/jfontan/log-server/process"
)

func main() {
	p := process.LoadProcess("process.yml")
	p.Start()
	p.Wait()
}
