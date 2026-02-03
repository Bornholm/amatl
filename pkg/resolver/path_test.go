package resolver

import (
	"runtime"
	"testing"
)

func TestPath_Scheme(t *testing.T) {
	tests := []struct {
		name     string
		path     Path
		expected string
	}{
		{"HTTP URL", Path("http://example.com/path"), "http"},
		{"HTTPS URL", Path("https://example.com/path"), "https"},
		{"File URL", Path("file:///path/to/file"), "file"},
		{"FTP URL", Path("ftp://example.com/path"), "ftp"},
		{"Custom scheme", Path("custom://resource/path"), "custom"},
		{"Memory scheme", Path("memory://cache/key"), "memory"},
		{"Database scheme", Path("db://localhost/table"), "db"},
		{"Unix absolute path", Path("/path/to/file"), ""},
		{"Unix relative path", Path("./path/to/file"), ""},
		{"Windows absolute path", Path("C:\\path\\to\\file"), ""},
		{"Windows relative path", Path(".\\path\\to\\file"), ""},
		{"Empty path", Path(""), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.path.Scheme()
			if result != tt.expected {
				t.Errorf("Scheme() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPath_IsURL(t *testing.T) {
	tests := []struct {
		name     string
		path     Path
		expected bool
	}{
		{"HTTP URL", Path("http://example.com/path"), true},
		{"HTTPS URL", Path("https://example.com/path"), true},
		{"File URL", Path("file:///path/to/file"), true},
		{"FTP URL", Path("ftp://example.com/path"), true},
		{"Custom scheme", Path("custom://resource/path"), true},
		{"Memory scheme", Path("memory://cache/key"), true},
		{"Database scheme", Path("db://localhost/table"), true},
		{"S3 scheme", Path("s3://bucket/key"), true},
		{"Redis scheme", Path("redis://localhost:6379/0"), true},
		{"Unix absolute path", Path("/path/to/file"), false},
		{"Unix relative path", Path("./path/to/file"), false},
		{"Windows absolute path", Path("C:\\path\\to\\file"), false},
		{"Windows relative path", Path(".\\path\\to\\file"), false},
		{"Empty path", Path(""), false},
		{"Invalid URL", Path("not-a-url"), false},
		{"Single char scheme (invalid)", Path("c://path"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.path.IsURL()
			if result != tt.expected {
				t.Errorf("IsURL() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPath_IsAbs(t *testing.T) {
	tests := []struct {
		name     string
		path     Path
		expected bool
	}{
		{"HTTP URL", Path("http://example.com/path"), true},
		{"HTTPS URL", Path("https://example.com/path"), true},
		{"File URL", Path("file:///path/to/file"), true},
		{"Unix absolute path", Path("/path/to/file"), runtime.GOOS != "windows"},
		{"Unix relative path", Path("./path/to/file"), false},
		{"Unix relative path 2", Path("path/to/file"), false},
		{"Windows absolute path", Path("C:\\path\\to\\file"), runtime.GOOS == "windows"},
		{"Windows relative path", Path(".\\path\\to\\file"), false},
		{"Windows relative path 2", Path("path\\to\\file"), false},
		{"Empty path", Path(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.path.IsAbs()
			if result != tt.expected {
				t.Errorf("IsAbs() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPath_Host(t *testing.T) {
	tests := []struct {
		name     string
		path     Path
		expected string
	}{
		{"HTTP URL", Path("http://example.com/path"), "example.com"},
		{"HTTPS URL with port", Path("https://example.com:8080/path"), "example.com:8080"},
		{"File URL", Path("file:///path/to/file"), ""},
		{"File URL with host", Path("file://localhost/path/to/file"), "localhost"},
		{"Unix path", Path("/path/to/file"), ""},
		{"Windows path", Path("C:\\path\\to\\file"), ""},
		{"Empty path", Path(""), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.path.Host()
			if result != tt.expected {
				t.Errorf("Host() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPath_URLPath(t *testing.T) {
	tests := []struct {
		name     string
		path     Path
		expected string
	}{
		{"HTTP URL", Path("http://example.com/path/to/file"), "/path/to/file"},
		{"HTTPS URL", Path("https://example.com/path/to/file"), "/path/to/file"},
		{"File URL", Path("file:///path/to/file"), "/path/to/file"},
		{"Unix path", Path("/path/to/file"), "/path/to/file"},
		{"Windows path", Path("C:\\path\\to\\file"), "C:\\path\\to\\file"},
		{"Relative path", Path("./path/to/file"), "./path/to/file"},
		{"Empty path", Path(""), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.path.URLPath()
			if result != tt.expected {
				t.Errorf("URLPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPath_WithAuth(t *testing.T) {
	tests := []struct {
		name     string
		path     Path
		username string
		password string
		expected string
	}{
		{
			"HTTP URL with auth",
			Path("http://example.com/path"),
			"user",
			"pass",
			"http://user:pass@example.com/path",
		},
		{
			"HTTPS URL with auth",
			Path("https://example.com/path"),
			"user",
			"pass",
			"https://user:pass@example.com/path",
		},
		{
			"HTTP URL with username only",
			Path("http://example.com/path"),
			"user",
			"",
			"http://user:@example.com/path",
		},
		{
			"Unix path (no change)",
			Path("/path/to/file"),
			"user",
			"pass",
			"/path/to/file",
		},
		{
			"Windows path (no change)",
			Path("C:\\path\\to\\file"),
			"user",
			"pass",
			"C:\\path\\to\\file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.path.WithAuth(tt.username, tt.password)
			if result.String() != tt.expected {
				t.Errorf("WithAuth() = %v, want %v", result.String(), tt.expected)
			}
		})
	}
}

func TestPath_Dir(t *testing.T) {
	tests := []struct {
		name     string
		path     Path
		expected string
	}{
		{"HTTP URL", Path("http://example.com/path/to/file.txt"), "http://example.com/path/to"},
		{"HTTPS URL", Path("https://example.com/path/to/file.txt"), "https://example.com/path/to"},
		{"File URL", Path("file:///path/to/file.txt"), "file:///path/to"},
		{"Current dir file", Path("./file.txt"), "."},
		{"Simple file", Path("file.txt"), "."},
	}

	// Platform-specific tests
	if runtime.GOOS == "windows" {
		windowsTests := []struct {
			name     string
			path     Path
			expected string
		}{
			{"Unix absolute path on Windows", Path("/path/to/file.txt"), "\\path\\to"},
			{"Unix relative path on Windows", Path("./path/to/file.txt"), "path\\to"},
			{"Unix relative path 2 on Windows", Path("path/to/file.txt"), "path\\to"},
			{"Root path on Windows", Path("/file.txt"), "\\"},
			{"Windows absolute path", Path("C:\\path\\to\\file.txt"), "C:\\path\\to"},
			{"Windows relative path", Path(".\\path\\to\\file.txt"), "path\\to"},
			{"Windows root", Path("C:\\file.txt"), "C:\\"},
		}
		tests = append(tests, windowsTests...)
	} else {
		unixTests := []struct {
			name     string
			path     Path
			expected string
		}{
			{"Unix absolute path", Path("/path/to/file.txt"), "/path/to"},
			{"Unix relative path", Path("./path/to/file.txt"), "path/to"},
			{"Unix relative path 2", Path("path/to/file.txt"), "path/to"},
			{"Root path", Path("/file.txt"), "/"},
		}
		tests = append(tests, unixTests...)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.path.Dir()
			if result.String() != tt.expected {
				t.Errorf("Dir() = %v, want %v", result.String(), tt.expected)
			}
		})
	}
}

func TestPath_JoinPath(t *testing.T) {
	tests := []struct {
		name     string
		base     Path
		paths    []string
		expected string
	}{
		{
			"HTTP URL join",
			Path("http://example.com/base"),
			[]string{"path", "to", "file.txt"},
			"http://example.com/base/path/to/file.txt",
		},
		{
			"HTTPS URL join",
			Path("https://example.com/base"),
			[]string{"path", "to", "file.txt"},
			"https://example.com/base/path/to/file.txt",
		},
		{
			"File URL join",
			Path("file:///base"),
			[]string{"path", "to", "file.txt"},
			"file:///base/path/to/file.txt",
		},
	}

	// Platform-specific tests
	if runtime.GOOS == "windows" {
		windowsTests := []struct {
			name     string
			base     Path
			paths    []string
			expected string
		}{
			{
				"Unix absolute path join on Windows",
				Path("/base"),
				[]string{"path", "to", "file.txt"},
				"\\base\\path\\to\\file.txt",
			},
			{
				"Unix relative path join on Windows",
				Path("./base"),
				[]string{"path", "to", "file.txt"},
				"base\\path\\to\\file.txt",
			},
			{
				"Single path join on Windows",
				Path("/base"),
				[]string{"file.txt"},
				"\\base\\file.txt",
			},
			{
				"Empty paths on Windows",
				Path("/base"),
				[]string{},
				"\\base",
			},
			{
				"Windows absolute path join",
				Path("C:\\base"),
				[]string{"path", "to", "file.txt"},
				"C:\\base\\path\\to\\file.txt",
			},
			{
				"Windows relative path join",
				Path(".\\base"),
				[]string{"path", "to", "file.txt"},
				"base\\path\\to\\file.txt",
			},
		}
		tests = append(tests, windowsTests...)
	} else {
		unixTests := []struct {
			name     string
			base     Path
			paths    []string
			expected string
		}{
			{
				"Unix absolute path join",
				Path("/base"),
				[]string{"path", "to", "file.txt"},
				"/base/path/to/file.txt",
			},
			{
				"Unix relative path join",
				Path("./base"),
				[]string{"path", "to", "file.txt"},
				"base/path/to/file.txt",
			},
			{
				"Single path join",
				Path("/base"),
				[]string{"file.txt"},
				"/base/file.txt",
			},
			{
				"Empty paths",
				Path("/base"),
				[]string{},
				"/base",
			},
		}
		tests = append(tests, unixTests...)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.base.JoinPath(tt.paths...)
			if result.String() != tt.expected {
				t.Errorf("JoinPath() = %v, want %v", result.String(), tt.expected)
			}
		})
	}
}

func TestPath_Join(t *testing.T) {
	tests := []struct {
		name     string
		base     Path
		paths    []Path
		expected string
	}{
		{
			"HTTP URL join with Path",
			Path("http://example.com/base"),
			[]Path{Path("path"), Path("to"), Path("file.txt")},
			"http://example.com/base/path/to/file.txt",
		},
	}

	// Platform-specific tests
	if runtime.GOOS == "windows" {
		windowsTests := []struct {
			name     string
			base     Path
			paths    []Path
			expected string
		}{
			{
				"Unix path join with Path on Windows",
				Path("/base"),
				[]Path{Path("path"), Path("to"), Path("file.txt")},
				"\\base\\path\\to\\file.txt",
			},
			{
				"Relative path join with Path on Windows",
				Path("./base"),
				[]Path{Path("path"), Path("file.txt")},
				"base\\path\\file.txt",
			},
			{
				"Windows path join with Path",
				Path("C:\\base"),
				[]Path{Path("path"), Path("file.txt")},
				"C:\\base\\path\\file.txt",
			},
		}
		tests = append(tests, windowsTests...)
	} else {
		unixTests := []struct {
			name     string
			base     Path
			paths    []Path
			expected string
		}{
			{
				"Unix path join with Path",
				Path("/base"),
				[]Path{Path("path"), Path("to"), Path("file.txt")},
				"/base/path/to/file.txt",
			},
			{
				"Relative path join with Path",
				Path("./base"),
				[]Path{Path("path"), Path("file.txt")},
				"base/path/file.txt",
			},
		}
		tests = append(tests, unixTests...)
	}

	// Add Windows-specific test
	if runtime.GOOS == "windows" {
		windowsTest := struct {
			name     string
			base     Path
			paths    []Path
			expected string
		}{
			"Windows path join with Path",
			Path("C:\\base"),
			[]Path{Path("path"), Path("file.txt")},
			"C:\\base\\path\\file.txt",
		}
		tests = append(tests, windowsTest)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.base.Join(tt.paths...)
			if result.String() != tt.expected {
				t.Errorf("Join() = %v, want %v", result.String(), tt.expected)
			}
		})
	}
}

func TestPath_CrossPlatformBehavior(t *testing.T) {
	// Test that URLs behave consistently across platforms
	urlTests := []struct {
		name string
		path Path
	}{
		{"HTTP URL", Path("http://example.com/path/to/file")},
		{"HTTPS URL", Path("https://example.com/path/to/file")},
		{"File URL", Path("file:///path/to/file")},
	}

	for _, tt := range urlTests {
		t.Run(tt.name, func(t *testing.T) {
			// URLs should always be absolute
			if !tt.path.IsAbs() {
				t.Errorf("URL %v should be absolute", tt.path)
			}

			// URLs should always be URLs
			if !tt.path.IsURL() {
				t.Errorf("URL %v should be detected as URL", tt.path)
			}

			// URLs should have consistent scheme detection
			scheme := tt.path.Scheme()
			if scheme == "" {
				t.Errorf("URL %v should have a scheme", tt.path)
			}
		})
	}
}

func TestPath_EdgeCases(t *testing.T) {
	tests := []struct {
		name string
		path Path
		test func(t *testing.T, p Path)
	}{
		{
			"Empty path",
			Path(""),
			func(t *testing.T, p Path) {
				if p.IsURL() {
					t.Error("Empty path should not be URL")
				}
				if p.IsAbs() {
					t.Error("Empty path should not be absolute")
				}
				if p.Scheme() != "" {
					t.Error("Empty path should have empty scheme")
				}
			},
		},
		{
			"Root path",
			Path("/"),
			func(t *testing.T, p Path) {
				if p.IsURL() {
					t.Error("Root path should not be URL")
				}
				// On Windows, Unix-style root path "/" is not considered absolute
				if runtime.GOOS != "windows" && !p.IsAbs() {
					t.Error("Root path should be absolute")
				}
				expectedDir := "/"
				if runtime.GOOS == "windows" {
					expectedDir = "\\"
				}
				if p.Dir().String() != expectedDir {
					t.Errorf("Root path dir should be %q, got %v", expectedDir, p.Dir())
				}
			},
		},
		{
			"Current directory",
			Path("."),
			func(t *testing.T, p Path) {
				if p.IsURL() {
					t.Error("Current dir should not be URL")
				}
				if p.IsAbs() {
					t.Error("Current dir should not be absolute")
				}
			},
		},
		{
			"URL with query and fragment",
			Path("https://example.com/path?query=value#fragment"),
			func(t *testing.T, p Path) {
				if !p.IsURL() {
					t.Error("URL with query should be URL")
				}
				if p.Scheme() != "https" {
					t.Errorf("Scheme should be https, got %v", p.Scheme())
				}
				if p.Host() != "example.com" {
					t.Errorf("Host should be example.com, got %v", p.Host())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t, tt.path)
		})
	}
}

// Benchmark tests for performance
func BenchmarkPath_Scheme(b *testing.B) {
	p := Path("https://example.com/path/to/file")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.Scheme()
	}
}

func BenchmarkPath_IsURL(b *testing.B) {
	p := Path("https://example.com/path/to/file")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.IsURL()
	}
}

func BenchmarkPath_JoinPath(b *testing.B) {
	p := Path("/base/path")
	paths := []string{"sub", "path", "file.txt"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.JoinPath(paths...)
	}
}

func BenchmarkPath_Dir(b *testing.B) {
	p := Path("/very/long/path/to/some/file.txt")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.Dir()
	}
}
