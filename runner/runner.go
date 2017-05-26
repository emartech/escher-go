package runner

import (
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
)

type Runner interface {
	Run()
}
type subProcess struct {
	port   string
	name   string
	args   []string
	signal chan os.Signal
}

func New(port string, name string, args []string, signal chan os.Signal) Runner {
	return &subProcess{port: "9292", name: name, args: args, signal: signal}
}

func (sp *subProcess) Run() {

	cmd := exec.Command(sp.name, sp.args...)

	cmd.Env = sp.envForSubProcess()
	stdout, err := cmd.StdoutPipe()
	sp.checkError(err)
	stderr, err := cmd.StderrPipe()
	sp.checkError(err)
	stdin, err := cmd.StderrPipe()
	sp.checkError(err)

	err = cmd.Start()
	sp.checkError(err)

	cmd.Wait()

	go func() { cmd.Process.Signal(<-sp.signal) }()

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	go io.Copy(os.Stdin, stdin)

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
	return append(env, "PORT="+sp.port)
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
