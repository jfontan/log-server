# Log Server

Log processing pipeline configurable with yaml files.

The pipeline is composed by nodes that pass the event data down the hierarchy. There are three classes of nodes:

* **source**: Generates event data. This can be a node that reads from a file, listens on a socket or directly generates events from a local source. It's also the node that parses the data into a dictionary of strings that will be consumed by the other nodes.

* **action**: This node has a parent (input event) and children. It can be used to transform the data, use it to do some calculations or filter the events used down the tree.

* **sink**: This node receives input events but does not forward it to any other node. It should be used for termination actions like saving events to a file or sending and email.

## Events

Events are parsed into a map. For example, this line:

```
t=2017-11-14t23:10:03+0000 lvl=info msg="job started" module=borges worker=2 job=015fbccb-f1b2-2903-3aac-4eafc2196a7d caller=archiver.go:76
```

is parsed as:

```yaml
t:      2017-11-14t23:10:03+0000
lvl:    info
msg:    job started
module: borges
worker: 2
job:    015fbccb-f1b2-2903-3aac-4eafc2196a7d caller=archiver.go:76
```


## Nodes

* **FileSource** (*source*): Reads log lines from a file. It's only able to parse [log15](https://github.com/inconshreveable/log15) format. The input file is specified in `args`.
* **RegexpFilter** (*process*): It filters events using regular expressions. Its input is a map specifying the key that has to check and as value is the regular expression it has to match. The filter map is specified in `map`.
* **PrintKeySink** (*sink*): Prints an specific event key. The key is specified in `args`.
* **CounterAction** (*action*): Creates a new counter named `args` and increments it each time an event passes the node.
* **PrintCountersSink** (*sink*): Prints all the counters each time it receives an event. It does not have arguments.

Example:

```yaml
---
- name: reader
  type: FileSource
  args: logs
  children:
    - name: count all
      type: CounterAction
      args: total
      children:
      - name: filter siva error
        type: RegexpFilter
        map:
          lvl: "^eror$"
          error: "index read failed"
        children:
          - name: print root
            type: PrintKeySink
            args: "root"
          - name: print msg
            type: PrintKeySink
            args: "msg"
          - name: count siva errors
            type: CounterAction
            args: siva_errors
            children:
              - name: print counters
                type: PrintCountersSink
```


