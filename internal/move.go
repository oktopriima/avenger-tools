package internal

import (
	"os"
	"os/exec"
)

func (c *createInternal) move() error {
	err := exec.Command("mkdir", "-p", c.TargetRoot).Run()
	if err != nil {
		return err
	}

	err = os.Rename(c.TempDir, c.TargetDir)
	if err != nil {
		return err
	}

	return nil
}
