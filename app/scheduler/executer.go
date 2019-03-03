package scheduler

import (
	"os/exec"
)

/*
Run func(command string, args ...string) ([]byte, error)
*/
func Run(command string, envs []string, args ...string) ([]byte, error) {

	cmd := exec.Command(command, args...)
	cmd.Env = envs

	output, err := cmd.StdoutPipe()
	// output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	var out = make([]byte, 0, 1024)
	for {
		tmp := make([]byte, 128)
		n, err := output.Read(tmp)
		out = append(out, tmp[:n]...)
		if err != nil {
			break
		}
	}

	if err = cmd.Wait(); err != nil {
		return nil, err
	}

	return out, nil

}
