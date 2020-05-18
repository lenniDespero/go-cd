package target

import (
	"fmt"
	"testing"

	"github.com/maraino/testify/require"
	"github.com/pkg/errors"

	"github.com/lenniDespero/go-cd/internal/pkg/host"
	"github.com/lenniDespero/go-cd/internal/pkg/pipe"
)

func prepareTarget() Config {
	return Config{
		Type: "local",
		Host: host.Config{},
		Path: "newPath",
		Pipe: []pipe.Config{
			{Name: "some links", Type: "links", Args: []interface{}{map[string]interface{}{"from": "pathFrom", "to": "pathTo"}}},
			{Name: "some Cmd", Type: "command", Args: []interface{}{map[string]interface{}{"command": "ls", "args": []string{"ls"}}}},
		},
	}
}

func check(t *testing.T, target Config, e error) {
	err := target.CheckConfig()
	require.Error(t, err, e)
}

func TestTarget_CheckConfig(t *testing.T) {
	target := prepareTarget()
	err := target.CheckConfig()
	require.Nil(t, err)
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
	require.Error(t, err, newErr)
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
	target.Pipe = []pipe.Config{}
	check(t, target, NoPipesError)
}

func TestTarget_CheckConfigBadPipes(t *testing.T) {
	target := prepareTarget()
	target.Pipe[1].Type = "xz"
	err := target.CheckConfig()
	require.NotNil(t, err)
}
