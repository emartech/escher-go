package runner

import (
	"io"
	"log"
	"os"
	"os/exec"
)

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

	stdout, err := cmd.StdoutPipe()
	checkError(err)
	stderr, err := cmd.StderrPipe()
	checkError(err)
	stdin, err := cmd.StderrPipe()
	checkError(err)

	err = cmd.Start()
	checkError(err)

	defer cmd.Wait()

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	go io.Copy(os.Stdin, stdin)

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
