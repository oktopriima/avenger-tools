package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func (c *createInternal) renamePackage() error {
	var files []string

	_ = filepath.Walk(c.TempDir, func(path string, info os.FileInfo, err error) error {
		// Skip .git folder
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	fmt.Printf("\n")
	bar := progressbar.NewOptions(len(files),
		progressbar.OptionSetDescription("Replacing package path..."),
		progressbar.OptionShowCount(),
		progressbar.OptionSetPredictTime(false),
	)

	return filepath.Walk(c.TempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		_ = bar.Add(1)
		if strings.HasSuffix(path, ".exe") || strings.HasSuffix(path, ".so") {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content := string(b)
		old := "github.com/oktopriima/marvel"

		if strings.Contains(content, old) {
			content = strings.ReplaceAll(content, old, c.PackageName)
			return os.WriteFile(path, []byte(content), info.Mode())
		}
		return nil
	})
}
