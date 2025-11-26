package hugo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gohugoio/hugo/commands"
)

type Service struct {
}

func New() Service {
	return Service{}
}

// runCommand executes a shell command in the current directory
func (s Service) runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run %s: %w", name, err)
	}
	return nil
}

// cleanup removes build artifacts and temporary files
func (s Service) cleanup() error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// List of paths to remove
	pathsToRemove := []string{
		filepath.Join(cwd, "public"),
		filepath.Join(cwd, "resources"),
		filepath.Join(cwd, "go.sum"),
		filepath.Join(cwd, ".hugo_build.lock"),
		filepath.Join(cwd, "themes"),
	}

	// Remove each path if it exists
	for _, path := range pathsToRemove {
		if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove %s: %w", path, err)
		}
	}

	return nil
}

// Execute runs the pre-build commands and then executes Hugo
func (s Service) Execute(args ...string) {
	// Run cleanup
	fmt.Println("Running cleanup...")
	if err := s.cleanup(); err != nil {
		fmt.Fprintf(os.Stderr, "Cleanup failed: %v\n", err)
		os.Exit(1)
	}

	// Run hugo mod clean using embedded Hugo
	fmt.Println("Running hugo mod clean...")
	commands.Execute([]string{"mod", "clean"})

	// Run go mod tidy (this needs to be external as it's a Go command)
	fmt.Println("Running go mod tidy...")
	if err := s.runCommand("go", "mod", "tidy"); err != nil {
		fmt.Fprintf(os.Stderr, "Go mod tidy failed: %v\n", err)
		os.Exit(1)
	}

	// Run hugo mod tidy using embedded Hugo
	fmt.Println("Running hugo mod tidy...")
	commands.Execute([]string{"mod", "tidy"})

	// Run hugo mod get using embedded Hugo
	fmt.Println("Running hugo mod get...")
	commands.Execute([]string{"mod", "get"})

	// If no args provided, use default flags
	if len(args) == 0 {
		args = []string{"--templateMetrics", "--ignoreCache", "--logLevel", "info"}
	}

	// Execute Hugo with the provided or default arguments
	fmt.Printf("Running hugo with args: %v\n", args)
	commands.Execute(args)
}
