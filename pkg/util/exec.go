package util

import (
	"os"
	"os/exec"
)

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
