package urlx

import (
	"net/url"
	"testing"
)

func TestDirWindowsPath(t *testing.T) {
	// Test case for the Windows path bug
	windowsPath := `c:\test\misc\amatl\test.md`

	u, err := url.Parse(windowsPath)
	if err != nil {
		t.Fatalf("Failed to parse Windows path: %v", err)
	}

	// Verify the URL is parsed as expected (opaque format)
	if u.Scheme != "c" {
		t.Errorf("Expected scheme 'c', got '%s'", u.Scheme)
	}

	if u.Opaque != `\test\misc\amatl\test.md` {
		t.Errorf("Expected opaque path, got '%s'", u.Opaque)
	}

	// Test the Dir function
	dirURL, err := Dir(u)
	if err != nil {
		t.Fatalf("Dir function failed: %v", err)
	}

	// The directory should be correctly extracted
	expectedDir := `\test\misc\amatl`
	if dirURL.Opaque != expectedDir {
		t.Errorf("Expected directory '%s', got '%s'", expectedDir, dirURL.Opaque)
	}

	// Verify the scheme is preserved
	if dirURL.Scheme != "c" {
		t.Errorf("Expected scheme 'c', got '%s'", dirURL.Scheme)
	}
}
