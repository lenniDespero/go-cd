package pipe

import (
	"log"
	"testing"

	"github.com/maraino/testify/require"
)

func preparePipe() Config {
	return Config{
		Name: "Some Config",
		Type: "links",
		Args: []interface{}{"some args"},
	}
}

func TestPipe_CheckConfig(t *testing.T) {
	pipe := preparePipe()
	log.Println(pipe)
	err := pipe.CheckConfig()
	require.Nil(t, err)
}

func TestPipe_CheckConfigNoName(t *testing.T) {
	pipe := preparePipe()
	pipe.Name = ""
	err := pipe.CheckConfig()
	require.Error(t, err, ErrNoPipeName)
}

func TestPipe_CheckConfigNoType(t *testing.T) {
	pipe := preparePipe()
	pipe.Type = ""
	err := pipe.CheckConfig()
	require.Error(t, err, ErrNoPipeType)
}

func TestPipe_CheckConfigWrongType(t *testing.T) {
	pipe := preparePipe()
	pipe.Type = "strange"
	err := pipe.CheckConfig()
	require.Error(t, err, ErrNotInPipes)
}

func TestPipe_CheckConfigNoArgs(t *testing.T) {
	pipe := preparePipe()
	pipe.Args = nil
	err := pipe.CheckConfig()
	require.Error(t, err, ErrNoPipeArgs)
}
