package process

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Step interface {
	Start()
	Out() chan *Event
	Done() chan bool
}

type Node struct {
	Name     string
	Type     string
	Args     interface{}
	Map      map[string]string
	Children []Node
}

type Nodes []Node

type Process struct {
	Nodes Nodes
	Done  []chan bool
}

type CreateStep func(chan *Event, Node) Step

type NodeType struct {
	Name   string
	Type   string
	Create CreateStep
}

var nodeTypes = map[string]NodeType{
	"FileSource": {
		Name: "FileSource",
		Type: "source",
		Create: func(c chan *Event, n Node) Step {
			fileName, ok := n.Args.(string)
			if !ok {
				panic("Argument invalid")
			}

			fileSource := GenFileSource(fileName)
			source := NewSource(fileSource)
			return source
		},
	},
	"RegexpFilter": {
		Name: "RegexpFilter",
		Type: "action",
		Create: func(c chan *Event, n Node) Step {
			filterArgs := n.Map

			regexpFilter := GenRegexpFilter(filterArgs)
			filter := NewAction(c, regexpFilter)
			return filter
		},
	},
	"PrintKeySink": {
		Name: "PrintKeySink",
		Type: "sink",
		Create: func(c chan *Event, n Node) Step {
			key, ok := n.Args.(string)
			if !ok {
				panic("Argument invalid")
			}

			printKey := GenPrintKeySink(key)
			sink := NewSink(c, printKey)
			return sink
		},
	},
	"CounterAction": {
		Name: "CounterAction",
		Type: "action",
		Create: func(c chan *Event, n Node) Step {
			name, ok := n.Args.(string)
			if !ok {
				panic("Argument invalid")
			}

			counterAction := GenCounterAction(name)
			action := NewAction(c, counterAction)
			return action
		},
	},
	"PrintCountersSink": {
		Name: "PrintCountersSink",
		Type: "sink",
		Create: func(c chan *Event, n Node) Step {
			printCounter := GenPrintCounters()
			sink := NewSink(c, printCounter)
			return sink
		},
	},
}

func LoadProcess(fileName string) *Process {
	nodes := make(Nodes, 0)

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, &nodes)
	if err != nil {
		panic(err)
	}

	return &Process{
		Nodes: nodes,
		Done:  nil,
	}
}

func (p *Process) startNodes(parent_ Step, nodes Nodes) {
	parent := parent_

	if parent != nil && len(nodes) > 1 {
		parent = NewTee(parent.Out())
		parent.Start()
	}

	for _, node := range nodes {
		t, ok := nodeTypes[node.Type]
		if !ok {
			panic("Node type " + node.Type + " not defined")
		}

		var n Step
		if parent == nil {
			n = t.Create(nil, node)
		} else {
			n = t.Create(parent.Out(), node)
		}

		n.Start()

		if t.Type == "sink" {
			p.Done = append(p.Done, n.Done())
		}

		if node.Children != nil {
			p.startNodes(n, node.Children)
		}
	}
}

func (p *Process) Start() (done []chan bool) {
	p.Done = make([]chan bool, 0)

	p.startNodes(nil, p.Nodes)

	return p.Done
}

func (p *Process) Wait() {
	println(len(p.Nodes))
	if p.Done == nil {
		return
	}

	for _, d := range p.Done {
		<-d
	}
}
