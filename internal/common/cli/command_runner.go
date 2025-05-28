package cli

import "os/exec"

type CommandRunner interface {
	Run(name string, args ...string) ([]byte, error)
}

type CliCommandRunner struct {
	exec func(name string, arg ...string) ([]byte, error)
}

func (c CliCommandRunner) Run(name string, args ...string) ([]byte, error) {
	return c.exec(name, args...)
}

func defaultExec(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

func NewCommandRunner() CommandRunner {
	return CliCommandRunner{
		exec: defaultExec,
	}
}
