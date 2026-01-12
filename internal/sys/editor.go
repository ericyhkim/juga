package sys

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// OpenEditor opens a temporary file with the given content in the user's preferred editor.
// It waits for the editor to close, then returns the updated content of the file.
func OpenEditor(initialContent string) (string, error) {
	tmpFile, err := os.CreateTemp("", "juga-edit-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	
	defer func() {
		tmpFile.Close()
		os.Remove(tmpPath)
	}()

	if _, err := tmpFile.WriteString(initialContent); err != nil {
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}
	tmpFile.Close() // Close explicitly so the editor can open it

	editor, err := getEditorCommand()
	if err != nil {
		return "", err
	}

	parts := strings.Fields(editor)
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid editor command")
	}

	path := parts[0]
	args := append(parts[1:], tmpPath)

	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor execution failed: %w", err)
	}

	content, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("failed to read back temp file: %w", err)
	}

	return string(content), nil
}

// getEditorCommand determines the best editor to use based on environment and OS.
func getEditorCommand() (string, error) {
	if visual := os.Getenv("VISUAL"); visual != "" {
		return visual, nil
	}

	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor, nil
	}

	candidates := []string{}
	if runtime.GOOS == "windows" {
		candidates = []string{"notepad"}
	} else {
		candidates = []string{"nano", "vi", "vim"}
	}

	for _, c := range candidates {
		if path, err := exec.LookPath(c); err == nil && path != "" {
			return c, nil
		}
	}

	return "", fmt.Errorf("no suitable editor found. Please set $EDITOR environment variable")
}
