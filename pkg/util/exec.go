package util

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ExecMust executes the given command and returns the output.
// panics if there is an error.
func ExecMustWithTrim(cmd string, args ...string) string {
	return strings.TrimSpace(ExecMust(cmd, args...))
}

// ExecMust executes the given command and returns the output.
// panics if there is an error.
func ExecMust(cmd string, args ...string) string {
	out, err := Exec(cmd, args...)

	if err != nil {
		panic(fmt.Sprintf("error while executing command: %s %s: %s", cmd, strings.Join(args, " "), err))
	}

	return out
}

// Exec executes the given command and returns the output.
// it shouldn't be used for vast outputs.
func Exec(cmd string, args ...string) (string, error) {
	var out, err = exec.Command(cmd, args...).Output()

	if err != nil {
		return "", err
	}

	return string(out), nil
}

// ExecStandard executes the given command and prints the output to stdout and stderr.
func ExecStandard(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}
