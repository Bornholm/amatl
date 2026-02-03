package resolver

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestRegistry_ResolveWithDifferentPaths(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "test content"

	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create registry with file resolver
	registry := NewRegistry()
	registry.Register("", &mockFileResolver{})
	registry.Register("file", &mockFileResolver{})
	registry.SetDefault("")

	ctx := context.Background()

	tests := []struct {
		name        string
		path        string
		shouldWork  bool
		description string
	}{
		{
			"Unix absolute path",
			testFile,
			true,
			"Should resolve absolute Unix path",
		},
		{
			"Unix relative path",
			"./test.txt",
			false, // Will fail because file doesn't exist in current dir
			"Relative path without working directory",
		},
		{
			"File URL absolute",
			"file://" + testFile,
			runtime.GOOS != "windows", // Skip on Windows due to file URL complexity
			"Should resolve file:// URL with absolute path",
		},
		{
			"HTTP URL",
			"http://example.com/test.txt",
			false, // Will fail because we don't have HTTP resolver
			"HTTP URL without HTTP resolver",
		},
		{
			"HTTPS URL",
			"https://example.com/test.txt",
			false, // Will fail because we don't have HTTPS resolver
			"HTTPS URL without HTTPS resolver",
		},
	}

	// Add Windows-specific tests
	if runtime.GOOS == "windows" {
		windowsTests := []struct {
			name        string
			path        string
			shouldWork  bool
			description string
		}{
			{
				"Windows absolute path",
				testFile, // This will be a Windows path on Windows
				true,
				"Should resolve absolute Windows path",
			},
			{
				"Windows relative path",
				".\\test.txt",
				false,
				"Windows relative path without working directory",
			},
		}
		tests = append(tests, windowsTests...)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := Path(tt.path)
			reader, err := registry.Resolve(ctx, path)

			if tt.shouldWork {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
					return
				}
				defer reader.Close()

				content, err := io.ReadAll(reader)
				if err != nil {
					t.Errorf("Failed to read content: %v", err)
					return
				}

				if string(content) != testContent {
					t.Errorf("Expected content %q, got %q", testContent, string(content))
				}
			} else {
				if err == nil {
					reader.Close()
					t.Errorf("Expected error but got success")
				}
			}
		})
	}
}

func TestRegistry_ResolveWithWorkingDirectory(t *testing.T) {
	// Create a temporary directory structure
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	testFile := filepath.Join(subDir, "test.txt")
	testContent := "test content with working dir"

	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create registry with file resolver
	registry := NewRegistry()
	registry.Register("", &mockFileResolver{})
	registry.SetDefault("")

	// Test with working directory
	workDir := Path(subDir)
	ctx := WithWorkDir(context.Background(), workDir)

	tests := []struct {
		name        string
		path        string
		shouldWork  bool
		description string
	}{
		{
			"Relative path with working directory",
			"test.txt",
			true,
			"Should resolve relative path with working directory",
		},
		{
			"Relative path with ./ prefix",
			"./test.txt",
			true,
			"Should resolve ./relative path with working directory",
		},
		{
			"Absolute path ignores working directory",
			testFile,
			true,
			"Absolute path should work regardless of working directory",
		},
		{
			"Non-existent relative file",
			"nonexistent.txt",
			false,
			"Should fail for non-existent relative file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := Path(tt.path)
			reader, err := registry.Resolve(ctx, path)

			if tt.shouldWork {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
					return
				}
				defer reader.Close()

				content, err := io.ReadAll(reader)
				if err != nil {
					t.Errorf("Failed to read content: %v", err)
					return
				}

				if string(content) != testContent {
					t.Errorf("Expected content %q, got %q", testContent, string(content))
				}
			} else {
				if err == nil {
					reader.Close()
					t.Errorf("Expected error but got success")
				}
			}
		})
	}
}

func TestRegistry_SchemeResolution(t *testing.T) {
	registry := NewRegistry()

	// Register different resolvers for different schemes
	registry.Register("mock", &mockResolver{content: "mock content"})
	registry.Register("test", &mockResolver{content: "test content"})
	registry.Register("", &mockResolver{content: "default content"})
	registry.SetDefault("")

	ctx := context.Background()

	tests := []struct {
		name            string
		path            string
		expectedContent string
		shouldError     bool
		description     string
	}{
		{
			"Mock scheme",
			"mock://example.com/path",
			"mock content",
			false,
			"Should use mock resolver for mock:// URLs",
		},
		{
			"Test scheme",
			"test://example.com/path",
			"test content",
			false,
			"Should use test resolver for test:// URLs",
		},
		{
			"No scheme uses default",
			"/path/to/file",
			"default content",
			false,
			"Should use default resolver for paths without scheme",
		},
		{
			"Unknown scheme should error",
			"unknown://example.com/path",
			"",
			true,
			"Should error for unknown schemes when no default resolver can handle it",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := Path(tt.path)
			reader, err := registry.Resolve(ctx, path)

			if tt.shouldError {
				if err == nil {
					reader.Close()
					t.Errorf("Expected error but got success")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			defer reader.Close()

			content, err := io.ReadAll(reader)
			if err != nil {
				t.Errorf("Failed to read content: %v", err)
				return
			}

			if string(content) != tt.expectedContent {
				t.Errorf("Expected content %q, got %q", tt.expectedContent, string(content))
			}
		})
	}
}

func TestContextWorkDir(t *testing.T) {
	tests := []struct {
		name     string
		workDir  Path
		expected Path
	}{
		{"Unix absolute path", Path("/tmp/test"), Path("/tmp/test")},
		{"Unix relative path", Path("./test"), Path("./test")},
		{"Empty path", Path(""), Path("")},
		{"URL path", Path("file:///tmp/test"), Path("file:///tmp/test")},
	}

	// Add Windows test
	if runtime.GOOS == "windows" {
		tests = append(tests, struct {
			name     string
			workDir  Path
			expected Path
		}{"Windows absolute path", Path("C:\\temp\\test"), Path("C:\\temp\\test")})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := WithWorkDir(context.Background(), tt.workDir)
			result := ContextWorkDir(ctx)

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestContextWorkDir_NoWorkDir(t *testing.T) {
	ctx := context.Background()
	result := ContextWorkDir(ctx)

	if result != "" {
		t.Errorf("Expected empty path, got %q", result)
	}
}

func TestDefaultResolver(t *testing.T) {
	// Test that default resolver exists and can be extended
	if DefaultResolver == nil {
		t.Error("DefaultResolver should not be nil")
	}

	// Test extending the default resolver
	extended := DefaultResolver.Extend(
		func() (scheme string, resolver Resolver) {
			return "test", &mockResolver{content: "extended content"}
		},
	)

	if extended == nil {
		t.Error("Extended resolver should not be nil")
	}

	// Test that extended resolver works
	ctx := context.Background()
	path := Path("test://example.com/path")
	reader, err := extended.Resolve(ctx, path)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		t.Errorf("Failed to read content: %v", err)
		return
	}

	if string(content) != "extended content" {
		t.Errorf("Expected 'extended content', got %q", string(content))
	}
}

// Mock resolvers for testing

type mockResolver struct {
	content string
	err     error
}

func (m *mockResolver) Resolve(ctx context.Context, path Path) (io.ReadCloser, error) {
	if m.err != nil {
		return nil, m.err
	}
	return io.NopCloser(strings.NewReader(m.content)), nil
}

type mockFileResolver struct{}

func (m *mockFileResolver) Resolve(ctx context.Context, path Path) (io.ReadCloser, error) {
	// Get the actual file path, handling file:// URLs
	filePath := path.String()
	scheme := path.Scheme()

	if scheme == "file" {
		// For file:// URLs, we need to handle both absolute and relative paths
		if u, err := path.URL(); err == nil {
			if u.Host != "" && u.Host != "localhost" {
				// Handle file://host/path format (relative paths like file://testdata/test.txt)
				filePath = u.Host + u.Path
			} else {
				// Handle file:///path format (absolute paths)
				filePath = u.Path
				// On Windows, convert /C:/path to C:/path, then let filepath handle separators
				if len(filePath) > 3 && filePath[0] == '/' && len(filePath) > 2 && filePath[2] == ':' {
					filePath = filePath[1:] // Remove leading slash: /C:/path -> C:/path
				}
				// Convert forward slashes to backslashes on Windows
				filePath = filepath.FromSlash(filePath)
			}
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Benchmark tests
func BenchmarkRegistry_Resolve(b *testing.B) {
	registry := NewRegistry()
	registry.Register("mock", &mockResolver{content: "benchmark content"})
	registry.SetDefault("mock")

	ctx := context.Background()
	path := Path("mock://example.com/path")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader, err := registry.Resolve(ctx, path)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
		reader.Close()
	}
}

func BenchmarkContextWorkDir(b *testing.B) {
	workDir := Path("/tmp/test/benchmark")
	ctx := WithWorkDir(context.Background(), workDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ContextWorkDir(ctx)
	}
}

// Test concurrent access
func TestRegistry_ConcurrentAccess(t *testing.T) {
	registry := NewRegistry()
	registry.Register("mock", &mockResolver{content: "concurrent content"})
	registry.SetDefault("mock")

	ctx := context.Background()
	path := Path("mock://example.com/concurrent")

	// Run multiple goroutines concurrently
	const numGoroutines = 10
	const numIterations = 100

	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer func() { done <- true }()

			for j := 0; j < numIterations; j++ {
				reader, err := registry.Resolve(ctx, path)
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}

				content, err := io.ReadAll(reader)
				reader.Close()

				if err != nil {
					t.Errorf("Failed to read content: %v", err)
					return
				}

				if string(content) != "concurrent content" {
					t.Errorf("Expected 'concurrent content', got %q", string(content))
					return
				}
			}
		}()
	}

	// Wait for all goroutines to complete with timeout
	timeout := time.After(10 * time.Second)
	for i := 0; i < numGoroutines; i++ {
		select {
		case <-done:
			// Goroutine completed successfully
		case <-timeout:
			t.Fatal("Test timed out waiting for goroutines to complete")
		}
	}
}

func TestRegistry_CustomSchemes(t *testing.T) {
	registry := NewRegistry()

	// Register custom resolvers for various schemes
	registry.Register("memory", &mockResolver{content: "memory data"})
	registry.Register("cache", &mockResolver{content: "cached data"})
	registry.Register("db", &mockResolver{content: "database record"})
	registry.Register("s3", &mockResolver{content: "s3 object"})
	registry.Register("redis", &mockResolver{content: "redis value"})
	registry.Register("custom", &mockResolver{content: "custom resource"})

	ctx := context.Background()

	tests := []struct {
		name            string
		path            string
		expectedContent string
		description     string
	}{
		{
			"Memory scheme",
			"memory://cache/user:123",
			"memory data",
			"Should resolve memory:// URLs",
		},
		{
			"Cache scheme",
			"cache://local/session:abc",
			"cached data",
			"Should resolve cache:// URLs",
		},
		{
			"Database scheme",
			"db://localhost:5432/users/123",
			"database record",
			"Should resolve db:// URLs",
		},
		{
			"S3 scheme",
			"s3://my-bucket/path/to/object.json",
			"s3 object",
			"Should resolve s3:// URLs",
		},
		{
			"Redis scheme",
			"redis://localhost:6379/0/key",
			"redis value",
			"Should resolve redis:// URLs",
		},
		{
			"Custom scheme",
			"custom://service/resource/123",
			"custom resource",
			"Should resolve custom:// URLs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := Path(tt.path)

			// Verify the path is correctly identified as a URL
			if !path.IsURL() {
				t.Errorf("Path %q should be identified as URL", path)
			}

			// Verify scheme extraction
			scheme := path.Scheme()
			expectedScheme := strings.Split(tt.path, "://")[0]
			if scheme != expectedScheme {
				t.Errorf("Expected scheme %q, got %q", expectedScheme, scheme)
			}

			// Test resolution
			reader, err := registry.Resolve(ctx, path)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			defer reader.Close()

			content, err := io.ReadAll(reader)
			if err != nil {
				t.Errorf("Failed to read content: %v", err)
				return
			}

			if string(content) != tt.expectedContent {
				t.Errorf("Expected content %q, got %q", tt.expectedContent, string(content))
			}
		})
	}
}

func TestRegistry_CustomSchemeWithAuth(t *testing.T) {
	registry := NewRegistry()
	registry.Register("secure", &mockAuthResolver{})

	ctx := context.Background()

	tests := []struct {
		name        string
		path        string
		username    string
		password    string
		description string
	}{
		{
			"Custom scheme with auth",
			"secure://api.example.com/resource",
			"user",
			"pass",
			"Should handle authentication in custom schemes",
		},
		{
			"Custom scheme with complex auth",
			"secure://service.internal/data/123",
			"admin",
			"secret123",
			"Should handle complex authentication",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := Path(tt.path)
			authPath := path.WithAuth(tt.username, tt.password)

			// Verify auth was added
			if !strings.Contains(authPath.String(), tt.username) {
				t.Errorf("Auth path should contain username %q", tt.username)
			}

			reader, err := registry.Resolve(ctx, authPath)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			defer reader.Close()

			content, err := io.ReadAll(reader)
			if err != nil {
				t.Errorf("Failed to read content: %v", err)
				return
			}

			expectedContent := "authenticated as " + tt.username
			if string(content) != expectedContent {
				t.Errorf("Expected content %q, got %q", expectedContent, string(content))
			}
		})
	}
}

func TestRegistry_SchemeRegistrationAndOverride(t *testing.T) {
	registry := NewRegistry()

	// Register initial resolver
	registry.Register("test", &mockResolver{content: "original"})

	ctx := context.Background()
	path := Path("test://example.com/resource")

	// Test original resolver
	reader, err := registry.Resolve(ctx, path)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	content, _ := io.ReadAll(reader)
	reader.Close()

	if string(content) != "original" {
		t.Errorf("Expected 'original', got %q", string(content))
	}

	// Override with new resolver
	registry.Register("test", &mockResolver{content: "overridden"})

	// Test overridden resolver
	reader, err = registry.Resolve(ctx, path)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	content, _ = io.ReadAll(reader)
	reader.Close()

	if string(content) != "overridden" {
		t.Errorf("Expected 'overridden', got %q", string(content))
	}
}

// Mock resolver with authentication support
type mockAuthResolver struct{}

func (m *mockAuthResolver) Resolve(ctx context.Context, path Path) (io.ReadCloser, error) {
	// Extract username from URL
	if u, err := path.URL(); err == nil && u.User != nil {
		username := u.User.Username()
		content := "authenticated as " + username
		return io.NopCloser(strings.NewReader(content)), nil
	}

	return io.NopCloser(strings.NewReader("no authentication")), nil
}

func TestRegistry_SchemeResolutionAfterPathJoining(t *testing.T) {
	registry := NewRegistry()

	// Register HTTP and file resolvers
	registry.Register("https", &mockResolver{content: "https content"})
	registry.Register("http", &mockResolver{content: "http content"})
	registry.Register("", &mockResolver{content: "file content"})
	registry.SetDefault("")

	ctx := context.Background()

	tests := []struct {
		name            string
		workDir         string
		relativePath    string
		expectedContent string
		expectedScheme  string
		description     string
	}{
		{
			"Relative path with HTTPS working directory",
			"https://example.com/base/dir",
			"../resource.txt",
			"https content",
			"https",
			"Should resolve to HTTPS when joined with HTTPS working directory",
		},
		{
			"Relative path with HTTP working directory",
			"http://example.com/base/dir",
			"../resource.txt",
			"http content",
			"http",
			"Should resolve to HTTP when joined with HTTP working directory",
		},
		{
			"Relative path with file working directory",
			"/tmp/base/dir",
			"../resource.txt",
			"file content",
			"",
			"Should resolve to file when joined with file working directory",
		},
		{
			"Absolute HTTPS path ignores working directory",
			"/tmp/base/dir",
			"https://example.com/resource.txt",
			"https content",
			"https",
			"Absolute HTTPS path should use HTTPS resolver regardless of working directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set working directory
			workDir := Path(tt.workDir)
			ctxWithWorkDir := WithWorkDir(ctx, workDir)

			// Resolve relative path
			relativePath := Path(tt.relativePath)
			reader, err := registry.Resolve(ctxWithWorkDir, relativePath)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			defer reader.Close()

			content, err := io.ReadAll(reader)
			if err != nil {
				t.Errorf("Failed to read content: %v", err)
				return
			}

			if string(content) != tt.expectedContent {
				t.Errorf("Expected content %q, got %q", tt.expectedContent, string(content))
			}

			// Verify the resolved path has the expected scheme
			resolvedPath := relativePath
			if workDir != "" && !relativePath.IsAbs() {
				resolvedPath = workDir.JoinPath(relativePath.String())
			}

			actualScheme := resolvedPath.Scheme()
			if actualScheme != tt.expectedScheme {
				t.Errorf("Expected scheme %q, got %q for resolved path %q", tt.expectedScheme, actualScheme, resolvedPath)
			}
		})
	}
}
