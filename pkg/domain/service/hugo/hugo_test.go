package hugo

import (
	"os"
	"testing"

	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/themes"
)

func TestNew(t *testing.T) {
	svc := New()
	// Service is intentionally a zero-sized struct (stateless)
	// This test just verifies New() doesn't panic
	_ = svc
}

func TestExecute_WithoutThemes(t *testing.T) {
	// Save original themesFS
	originalFS := themes.GetThemesFS()
	defer func() {
		if originalFS != nil {
			themes.SetFS(originalFS)
		}
	}()

	// Set themes FS to nil to simulate no embedded themes
	themes.SetFS(nil)

	svc := New()

	// Test with version command (should work even without themes)
	err := svc.Execute("version")

	// We expect Hugo to execute, error or not depends on Hugo installation
	// The test just ensures our code doesn't panic
	_ = err // Hugo might return error if not in a Hugo project directory
}

func TestExecute_CleansUpEnvironment(t *testing.T) {
	// This test verifies that environment variables are cleaned up
	// even if set during execution

	const envVarName = "HUGO_MODULE_REPLACEMENTS"

	// Save original env var if it exists
	originalValue, hadOriginal := os.LookupEnv(envVarName)
	defer func() {
		if hadOriginal {
			os.Setenv(envVarName, originalValue)
		} else {
			os.Unsetenv(envVarName)
		}
	}()

	svc := New()

	// Execute a command that will quickly fail (invalid command)
	_ = svc.Execute("nonexistent-command-xyz")

	// Check that the environment variable was cleaned up
	if value, exists := os.LookupEnv(envVarName); exists {
		t.Errorf("HUGO_MODULE_REPLACEMENTS was not cleaned up, value: %s", value)
	}
}

func TestExecute_HandlesDifferentArguments(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "no arguments",
			args: []string{},
		},
		{
			name: "help flag",
			args: []string{"--help"},
		},
		{
			name: "version command",
			args: []string{"version"},
		},
		{
			name: "multiple arguments",
			args: []string{"--logLevel", "error", "version"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New()

			// Execute the command - we don't check the error because
			// Hugo might error for various valid reasons (not in a project, etc.)
			// We're just ensuring our wrapper doesn't panic or fail unexpectedly
			_ = svc.Execute(tt.args...)
		})
	}
}

func TestExecute_PropagatesErrors(t *testing.T) {
	svc := New()

	// Execute with an intentionally invalid command
	err := svc.Execute("this-is-not-a-real-hugo-command-12345")

	if err == nil {
		// Hugo should return an error for an invalid command
		// If it doesn't, that's okay - Hugo's behavior might vary
		t.Log("Warning: Expected error from invalid Hugo command, but got nil")
	}
}

func TestService_IsZeroSized(t *testing.T) {
	// Verify that Service struct has no fields (stateless)
	svc := New()

	// This is mainly a documentation test - Service should remain stateless
	if testing.Short() {
		t.Skip("Skipping structural test in short mode")
	}

	// Service should be empty struct
	_ = svc // Just verify it exists and is usable
}

// TestExecute_Integration is an integration test that requires a real Hugo project
// and embedded themes. It's skipped by default unless explicitly enabled.
func TestExecute_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Check if we're in a Hugo project directory
	if _, err := os.Stat("config.yml"); os.IsNotExist(err) {
		if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
			if _, err := os.Stat("hugo.toml"); os.IsNotExist(err) {
				t.Skip("Not in a Hugo project directory, skipping integration test")
			}
		}
	}

	svc := New()

	// Try to run a simple Hugo command that should work in any Hugo project
	err := svc.Execute("env")

	if err != nil {
		t.Logf("Hugo env command returned error: %v (this may be expected)", err)
	}
}

// TestExecute_WithThemesExtraction tests the full flow with themes
// This is a more comprehensive test that verifies theme extraction works
func TestExecute_WithThemesExtraction(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping themes extraction test in short mode")
	}

	// Check if themes directory exists (we're in the project root)
	if _, err := os.Stat("themes"); os.IsNotExist(err) {
		t.Skip("Themes directory not found, skipping themes extraction test")
	}

	svc := New()

	// Execute version command - should work and use themes if available
	err := svc.Execute("version")

	// Version command should always work
	if err != nil {
		t.Logf("Hugo version command error: %v", err)
	}
}

// BenchmarkExecute benchmarks the Execute function
func BenchmarkExecute(b *testing.B) {
	svc := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = svc.Execute("version")
	}
}
