package cli

import "os/exec"

type CommandRunner interface {
	Run(name string, args ...string) ([]byte, error)
}

type cliCommandRunner struct {
	exec func(name string, arg ...string) ([]byte, error)
}

func (c cliCommandRunner) Run(name string, args ...string) ([]byte, error) {
	return c.exec(name, args...)
}

func defaultExec(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

func NewCommandRunner() CommandRunner {
	return cliCommandRunner{
		exec: defaultExec,
	}
}
