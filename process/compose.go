package process

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Node struct {
	Name     string
	Type     string
	Args     interface{}
	Map      map[string]string
	Children []Node
}

type Nodes []Node

type CreateBasicNode func(chan *Event, Node) BasicNode

type NodeType struct {
	Name   string
	Type   string
	Create CreateBasicNode
}

var nodeTypes = map[string]NodeType{
	"FileSource": {
		Name: "FileSource",
		Type: "source",
		Create: func(c chan *Event, n Node) BasicNode {
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
		Type: "process",
		Create: func(c chan *Event, n Node) BasicNode {
			filterArgs := n.Map

			regexpFilter := GenRegexpFilter(filterArgs)
			filter := NewProcess(c, regexpFilter)
			return filter
		},
	},
	"PrintKeySink": {
		Name: "PrintKeySink",
		Type: "sink",
		Create: func(c chan *Event, n Node) BasicNode {
			key, ok := n.Args.(string)
			if !ok {
				panic("Argument invalid")
			}

			printKey := GenPrintKeySink(key)
			sink := NewSink(c, printKey)
			return sink
		},
	},
}

func LoadNodes(fileName string) Nodes {
	nodes := make(Nodes, 0)

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, &nodes)
	if err != nil {
		panic(err)
	}

	return nodes
}

func startNodes(parent_ BasicNode, nodes Nodes, done *[]chan bool) {
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

		var n BasicNode
		if parent == nil {
			n = t.Create(nil, node)
		} else {
			n = t.Create(parent.Out(), node)
		}

		n.Start()

		if t.Type == "sink" {
			*done = append(*done, n.Done())
		}

		if node.Children != nil {
			startNodes(n, node.Children, done)
		}
	}
}

func (n Nodes) Start() (done []chan bool) {
	d := make([]chan bool, 0)

	startNodes(nil, n, &d)

	return d
}
