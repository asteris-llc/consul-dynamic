package consul

import (
	"os/exec"
)

func RestartConsul() error {
	path, err := exec.LookPath("systemctl")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, "restart", "consul")
	return cmd.Run()
}
