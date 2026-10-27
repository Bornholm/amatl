package destination

import (
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	testCases := map[string]string{
		"./img.png":                  "./img.png",
		"./my%20image.png":           "./my image.png",
		`./image\(1\).png`:           "./image(1).png",
		"./image%20%281%29.png":      "./image (1).png",
		"./été.png":                  "./été.png",
		"./100%.png":                 "./100%.png",
		"https://example.com/a%20b":  "https://example.com/a%20b",
		"file:///tmp/my%20image.png": "file:///tmp/my%20image.png",
		"/abs/dir/my image.png":      "/abs/dir/my image.png",
	}

	for raw, expected := range testCases {
		if g := Decode(raw); g != expected {
			t.Errorf("Decode(%q): expected %q, got %q", raw, expected, g)
		}
	}
}

func TestEncodeRoundTrip(t *testing.T) {
	for _, path := range []string{
		"/abs/dir/img.png",
		"/abs/dir/my image.png",
		"/abs/dir/image (1).png",
		"/abs/dir/100%.png",
		"/abs/dir/a%20b.png",
		"/abs/dir/été.png",
		"/abs/<dir>/img.png",
	} {
		encoded := Encode(path)

		if strings.ContainsAny(encoded, " ()<>\\") {
			t.Errorf("Encode(%q): unexpected character in %q", path, encoded)
		}

		if g := Decode(encoded); g != path {
			t.Errorf("Decode(Encode(%q)): got %q via %q", path, g, encoded)
		}
	}
}
