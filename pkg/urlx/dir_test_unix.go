//go:build !windows
// +build !windows

package urlx

import (
	"net/url"
	"testing"
)

func TestDirUnixPath(t *testing.T) {
	// Test case for Unix paths (should still work)
	unixPath := `/home/user/documents/file.md`

	u, err := url.Parse(unixPath)
	if err != nil {
		t.Fatalf("Failed to parse Unix path: %v", err)
	}

	// Test the Dir function
	dirURL, err := Dir(u)
	if err != nil {
		t.Fatalf("Dir function failed: %v", err)
	}

	// The directory should be correctly extracted
	expectedDir := `/home/user/documents`
	if dirURL.Path != expectedDir {
		t.Errorf("Expected directory '%s', got '%s'", expectedDir, dirURL.Path)
	}
}

func TestDirFileURL(t *testing.T) {
	// Test case for file:// URLs (should still work)
	fileURL := `file:///home/user/documents/file.md`

	u, err := url.Parse(fileURL)
	if err != nil {
		t.Fatalf("Failed to parse file URL: %v", err)
	}

	// Test the Dir function
	dirURL, err := Dir(u)
	if err != nil {
		t.Fatalf("Dir function failed: %v", err)
	}

	// The directory should be correctly extracted
	expectedDir := `/home/user/documents`
	if dirURL.Path != expectedDir {
		t.Errorf("Expected directory '%s', got '%s'", expectedDir, dirURL.Path)
	}

	// Verify the scheme is preserved
	if dirURL.Scheme != "file" {
		t.Errorf("Expected scheme 'file', got '%s'", dirURL.Scheme)
	}
}
