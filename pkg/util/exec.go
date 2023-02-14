package util

import "os/exec"

// Exec executes the given command and returns the output.
// it shouldn't be used for vast outputs.
func Exec(cmd string, args ...string) (string, error) {
	var out, err = exec.Command(cmd, args...).Output()

	if err != nil {
		return "", err
	}

	return string(out), nil
}
