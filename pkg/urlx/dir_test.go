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
