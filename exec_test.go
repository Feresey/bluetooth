package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

const (
	envHelper   = "GO_WANT_HELPER_PROCESS"
	envResult   = "GO_EXEC_RESULT"
	envExitCode = "GO_EXEC_EXIT_CODE"
)

type trapExec struct {
	text   string
	code   int
	called int
}

var result = make(map[string]*trapExec)

func TestHelperProcess(t *testing.T) {
	if os.Getenv(envHelper) != "1" {
		return
	}

	fmt.Println(os.Getenv(envResult))

	code, err := strconv.Atoi(os.Getenv(envExitCode))
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if code != 0 {
		os.Exit(code)
	}

	//bluetoothctl говно по сути, и какого-то хера он яростно кеширует всё, что не падает в stdout
	select {}
}

func helperCommandContext(ctx context.Context, name string, args ...string) (cmd *exec.Cmd) {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, name)
	cs = append(cs, args...)

	if ctx != context.TODO() {
		cmd = exec.CommandContext(ctx, os.Args[0], cs...)
	} else {
		cmd = exec.Command(os.Args[0], cs...)
	}

	res := result[strings.Join(cs[2:], " ")]

	res.called++

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("%s=%d", envHelper, 1),
		fmt.Sprintf("%s=%d", envExitCode, res.code),
		fmt.Sprintf("%s=%s", envResult, res.text),
	)

	return cmd

}

func helperCommand(name string, args ...string) *exec.Cmd {
	return helperCommandContext(context.TODO(), name, args...)
}
