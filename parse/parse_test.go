package parse

import (
	"reflect"
	"testing"
)

var (
	text = `t=2017-11-14t23:10:03+0000 lvl=info msg="initializing pack process" module=borges file=repos.txt output=repositories caller=packer.go:34
t=2017-11-14t23:10:03+0000 lvl=info msg=starting module=borges worker=1 caller=worker.go:37
t=2017-11-14t23:10:03+0000 lvl=info msg=starting module=borges worker=0 caller=worker.go:37
t=2017-11-14t23:10:03+0000 lvl=info msg=starting module=borges worker=2 caller=worker.go:37
t=2017-11-14t23:10:03+0000 lvl=info msg=starting module=borges worker=3 caller=worker.go:37
t=2017-11-14t23:10:03+0000 lvl=info msg="job started" module=borges worker=1 job=015fbccb-f1b2-3278-d997-038f7156d8b7 caller=archiver.go:76
t=2017-11-14t23:10:03+0000 lvl=info msg="job started" module=borges worker=0 job=015fbccb-f1b2-aff2-5d04-868e93735dfd caller=archiver.go:76
t=2017-11-14t23:10:03+0000 lvl=info msg="job started" module=borges worker=2 job=015fbccb-f1b2-2903-3aac-4eafc2196a7d caller=archiver.go:76`

	examples = map[string]map[string]string{
		`t=2017-11-14t23:10:03+0000 lvl=info msg="initializing pack process" module=borges file=repos.txt output=repositories caller=packer.go:34`: {
			"t":      "2017-11-14t23:10:03+0000",
			"lvl":    "info",
			"msg":    "initializing pack process",
			"module": "borges",
			"file":   "repos.txt",
			"output": "repositories",
			"caller": "packer.go:34",
		},
		`t=2017-11-14t23:10:03+0000 lvl=info msg=starting module=borges worker=1 caller=worker.go:37`: {
			"t":      "2017-11-14t23:10:03+0000",
			"lvl":    "info",
			"msg":    "starting",
			"module": "borges",
			"worker": "1",
			"caller": "worker.go:37",
		},
	}
)

// mapEqual can be used to pinpoint what is different
func mapEqual(a, b map[string]string) bool {
	if a == nil && b == nil {
		println("nil nil")
		return true
	} else if a == nil || b == nil {
		println("one nil")
		return false
	}

	aLen := len(a)
	bLen := len(b)

	if aLen != bLen {
		println("len")
		return false
	}

	for k, v := range a {
		if bV, ok := b[k]; !ok || v != bV {
			println(k, v, bV)
			println("key differs", k, v, bV)
			return false
		}
	}

	return true
}

func TestExamples(t *testing.T) {
	for k, v := range examples {
		res := ParseLine(k)

		if !reflect.DeepEqual(v, res) {
			t.Errorf("Result does not match for %q.\n  Expected: %q\n  Got:      %q", k, v, res)
		}
	}
}
