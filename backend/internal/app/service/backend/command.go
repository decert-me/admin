package backend

import (
	"bufio"
	"fmt"
	"os/exec"
)

func execCommand(dir string, command string, args ...string) (stdoutRes []string, stderrRes []string, err error) {
	cmd := exec.Command(command, args...)
	fmt.Println(cmd.Args) //显示运行的命令
	cmd.Dir = dir
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return stdoutRes, stderrRes, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return stdoutRes, stderrRes, err
	}
	err = cmd.Start()
	if err != nil {
		return stdoutRes, stderrRes, err
	}
	oReader := bufio.NewReader(stdout)
	eReader := bufio.NewReader(stderr)
	for {
		line, err := oReader.ReadString('\n')
		if err != nil {
			break
		}
		stdoutRes = append(stdoutRes, line)
	}

	for {
		line, err := eReader.ReadString('\n')
		if err != nil {
			break
		}
		stderrRes = append(stderrRes, line)
	}
	cmd.Wait()
	return stdoutRes, stderrRes, err
}
