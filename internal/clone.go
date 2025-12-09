package internal

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/schollz/progressbar/v3"
)

func (c *createInternal) clone() error {
	bar := progressbar.NewOptions(500,
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetDescription("Cloning repository..."),
		progressbar.OptionSetPredictTime(false),
	)

	pw := &ProgressWriter{bar: bar}

	_, err := git.PlainClone(c.TempDir, false, &git.CloneOptions{
		URL:      SourceRepository,
		Progress: pw,
	})

	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	_ = bar.Finish()
	fmt.Printf("\n")

	return nil
}
