package util

import (
	"fmt"
	"os/exec"

	"github.com/emicklei/tre"
	"github.com/google/shlex"
)

func RunCommand(line string, dryrun bool) (string, error) {
	parts, err := shlex.Split(line)
	if err != nil {
		return "", tre.New(err, "shlex.Split", "line", line)
	}
	for _, each := range parts {
		fmt.Println(each)
	}
	if dryrun {
		return "", nil
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	data, err := cmd.CombinedOutput()
	return string(data), err
}
