package internal

import (
	"github.com/schollz/progressbar/v3"
)

type ProgressWriter struct {
	bar *progressbar.ProgressBar
}

func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	_ = pw.bar.Add(1)
	return len(p), nil
}
