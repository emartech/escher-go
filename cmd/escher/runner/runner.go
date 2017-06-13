package runner

import (
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
)

type Runner interface {
	Run() *exec.Cmd
}
type subProcess struct {
	port   string
	name   string
	args   []string
	signal chan os.Signal
}

func New(port string, name string, args []string, signal chan os.Signal) Runner {
	return &subProcess{port, name, args, signal}
}

func (sp *subProcess) Run() *exec.Cmd {

	cmd := exec.Command(sp.name, sp.args...)
	cmd.Args[2] = sp.port

	stdout, err := cmd.StdoutPipe()
	sp.checkError(err)
	stderr, err := cmd.StderrPipe()
	sp.checkError(err)

	go func() { cmd.Process.Signal(<-sp.signal) }()

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	err = cmd.Start()
	sp.checkError(err)

	return cmd

}

func (sp *subProcess) checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (sp *subProcess) envForSubProcess() []string {
	return sp.addNewPort(sp.removeOldPortFromEnv())
}

func (sp *subProcess) addNewPort(env []string) []string {
	env = append(env, "PORT="+sp.port)
	return env
}

func (sp *subProcess) removeOldPortFromEnv() []string {
	var newEnv []string
	for _, v := range os.Environ() {
		if match, _ := regexp.MatchString("PORT", v); match == false {
			newEnv = append(newEnv, v)
		}
	}
	return newEnv
}
