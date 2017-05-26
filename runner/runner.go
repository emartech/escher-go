package runner

import "os/exec"

type Runner interface {
	Run()
}
type subProcess struct {
	port string
	name string
	args []string
}

func New(port string, name string, args []string) Runner {
	return &subProcess{port, name, args}
}

func (sp *subProcess) Run() {
	cmd := exec.Command(sp.name, sp.args...)
	cmd.Start()
}
