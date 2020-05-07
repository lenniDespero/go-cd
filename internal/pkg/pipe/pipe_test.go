package pipe

import (
	"log"
	"testing"
)

func preparePipe() Pipe {
	return Pipe{
		Name: "Some Pipe",
		Type: "links",
		Args: []interface{}{"some args"},
	}
}

func TestPipe_CheckConfig(t *testing.T) {
	pipe := preparePipe()
	log.Println(pipe)
	err := pipe.CheckConfig()
	if err != nil {
		t.Errorf("Unexpected error: %s, expected nil", err.Error())
	}
}

func TestPipe_CheckConfigNoName(t *testing.T) {
	pipe := preparePipe()
	pipe.Name = ""
	err := pipe.CheckConfig()
	if err != nil {
		if err != NoPipeName {
			t.Errorf("Unexpected error: %s, expected: %s", err.Error(), NoPipeName.Error())
		}
	} else {
		t.Errorf("Expected error, get nil")
	}
}

func TestPipe_CheckConfigNoType(t *testing.T) {
	pipe := preparePipe()
	pipe.Type = ""
	err := pipe.CheckConfig()
	if err != nil {
		if err != NoPipeTypeError {
			t.Errorf("Unexpected error: %s, expected: %s", err.Error(), NoPipeTypeError.Error())
		}
	} else {
		t.Errorf("Expected error, get nil")
	}
}

func TestPipe_CheckConfigWrongType(t *testing.T) {
	pipe := preparePipe()
	pipe.Type = "strange"
	err := pipe.CheckConfig()
	if err != nil {
		if err != NotInPipesError {
			t.Errorf("Unexpected error: %s, expected: %s", err.Error(), NotInPipesError.Error())
		}
	} else {
		t.Errorf("Expected error, get nil")
	}
}

func TestPipe_CheckConfigNoArgs(t *testing.T) {
	pipe := preparePipe()
	pipe.Args = nil
	err := pipe.CheckConfig()
	if err != nil {
		if err != NoPipeArgs {
			t.Errorf("Unexpected error: %s, expected: %s", err.Error(), NoPipeArgs.Error())
		}
	} else {
		t.Errorf("Expected error, get nil")
	}
}
