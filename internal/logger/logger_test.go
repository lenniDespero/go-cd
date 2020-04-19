package logger

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestDebug(t *testing.T) {
	var str bytes.Buffer
	log.SetOutput(&str)
	log.SetFlags(0)
	Debug("test")
	res := strings.TrimSuffix(str.String(), "\n")
	exp := fmt.Sprintf(DebugColor, "test")
	if res != exp {
		t.Errorf("expected: %s get %s", exp, res)
	}
}

func TestError(t *testing.T) {
	var str bytes.Buffer
	log.SetOutput(&str)
	log.SetFlags(0)
	Error("test")
	res := strings.TrimSuffix(str.String(), "\n")
	exp := fmt.Sprintf(ErrorColor, "test")
	if res != exp {
		t.Errorf("expected: %s get %s", exp, res)
	}
}

func TestNotice(t *testing.T) {
	var str bytes.Buffer
	log.SetOutput(&str)
	log.SetFlags(0)
	Notice("test")
	res := strings.TrimSuffix(str.String(), "\n")
	exp := fmt.Sprintf(NoticeColor, "test")
	if res != exp {
		t.Errorf("expected: %s get %s", exp, res)
	}
}

func TestInfo(t *testing.T) {
	var str bytes.Buffer
	log.SetOutput(&str)
	log.SetFlags(0)
	Info("test")
	res := strings.TrimSuffix(str.String(), "\n")
	exp := fmt.Sprintf(InfoColor, "test")
	if res != exp {
		t.Errorf("expected: %s get %s", exp, res)
	}
}

func TestWarn(t *testing.T) {
	var str bytes.Buffer
	log.SetOutput(&str)
	log.SetFlags(0)
	Warn("test")
	res := strings.TrimSuffix(str.String(), "\n")
	exp := fmt.Sprintf(WarningColor, "test")
	if res != exp {
		t.Errorf("expected: %s get %s", exp, res)
	}
}
