package target

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"

	"github.com/lenniDespero/go-cd/internal/pkg/host"

	"github.com/lenniDespero/go-cd/internal/pkg/pipe"
)

func prepareTarget() Target {
	return Target{
		Type: "local",
		Host: host.Host{},
		Path: "newPath",
		Pipe: []pipe.Pipe{
			{Name: "some links", Type: "links", Args: []interface{}{map[string]interface{}{"from": "pathFrom", "to": "pathTo"}}},
			{Name: "some Cmd", Type: "command", Args: []interface{}{map[string]interface{}{"command": "ls", "args": []string{"ls"}}}},
		},
	}
}

func check(t *testing.T, target Target, e error) {
	err := target.CheckConfig()
	if err != nil {
		if err != e {
			t.Errorf("Unexpected error: %s, expected: %s", err.Error(), e.Error())
		}
	} else {
		t.Errorf("Expected error, get nil")
	}
}

func TestTarget_CheckConfig(t *testing.T) {
	target := prepareTarget()
	err := target.CheckConfig()
	if err != nil {
		t.Errorf("Err: %s", err.Error())
	}
}

func TestTarget_CheckConfigNoType(t *testing.T) {
	target := prepareTarget()
	target.Type = ""
	check(t, target, NoTypeError)
}

func TestTarget_CheckConfigNotInTypes(t *testing.T) {
	target := prepareTarget()
	target.Type = "strange"
	newErr := errors.Wrap(NotInTypesError, fmt.Sprintf("type %s ", target.Type))
	err := target.CheckConfig()
	if err != nil {
		if err.Error() != newErr.Error() {
			t.Errorf("Errors  %s - %s not equal", err.Error(), newErr.Error())
		}
	}
}

func TestTarget_CheckConfigNoHost(t *testing.T) {
	target := prepareTarget()
	target.Type = "host"
	check(t, target, NoHostError)
}

func TestTarget_CheckConfigNoPath(t *testing.T) {
	target := prepareTarget()
	target.Path = ""
	check(t, target, NoPathError)
}

func TestTarget_CheckConfigNoPipes(t *testing.T) {
	target := prepareTarget()
	target.Pipe = []pipe.Pipe{}
	check(t, target, NoPipesError)
}

func TestTarget_CheckConfigBadPipes(t *testing.T) {
	target := prepareTarget()
	target.Pipe[1].Type = "xz"
	err := target.CheckConfig()
	if err == nil {
		t.Errorf("Expected err get nill")
	}
}
