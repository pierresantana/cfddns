package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFile(t *testing.T) {
	t.Run("loads variables", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, ".env")
		os.WriteFile(path, []byte("FOO_TEST_VAR=hello\nBAR_TEST_VAR=world\n"), 0644)

		os.Unsetenv("FOO_TEST_VAR")
		os.Unsetenv("BAR_TEST_VAR")
		defer os.Unsetenv("FOO_TEST_VAR")
		defer os.Unsetenv("BAR_TEST_VAR")

		if err := LoadFile(path); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := os.Getenv("FOO_TEST_VAR"); got != "hello" {
			t.Errorf("FOO_TEST_VAR = %q, want %q", got, "hello")
		}
		if got := os.Getenv("BAR_TEST_VAR"); got != "world" {
			t.Errorf("BAR_TEST_VAR = %q, want %q", got, "world")
		}
	})

	t.Run("does not override existing vars", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, ".env")
		os.WriteFile(path, []byte("EXISTING_TEST_VAR=new\n"), 0644)

		os.Setenv("EXISTING_TEST_VAR", "original")
		defer os.Unsetenv("EXISTING_TEST_VAR")

		if err := LoadFile(path); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := os.Getenv("EXISTING_TEST_VAR"); got != "original" {
			t.Errorf("EXISTING_TEST_VAR = %q, want %q", got, "original")
		}
	})

	t.Run("ignores missing file", func(t *testing.T) {
		err := LoadFile("/nonexistent/path/.env")
		if err != nil {
			t.Fatalf("expected nil error for missing file, got: %v", err)
		}
	})

	t.Run("ignores comments and blank lines", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, ".env")
		os.WriteFile(path, []byte("# comment\n\nVALID_TEST_VAR=yes\n"), 0644)

		os.Unsetenv("VALID_TEST_VAR")
		defer os.Unsetenv("VALID_TEST_VAR")

		if err := LoadFile(path); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := os.Getenv("VALID_TEST_VAR"); got != "yes" {
			t.Errorf("VALID_TEST_VAR = %q, want %q", got, "yes")
		}
	})

	t.Run("strips quotes from values", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, ".env")
		os.WriteFile(path, []byte(`QUOTED_TEST_VAR="quoted value"`+"\n"), 0644)

		os.Unsetenv("QUOTED_TEST_VAR")
		defer os.Unsetenv("QUOTED_TEST_VAR")

		if err := LoadFile(path); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := os.Getenv("QUOTED_TEST_VAR"); got != "quoted value" {
			t.Errorf("QUOTED_TEST_VAR = %q, want %q", got, "quoted value")
		}
	})
}
