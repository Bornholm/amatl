package node

import "testing"

func TestFormatDestination(t *testing.T) {
	testCases := map[string]string{
		"./img.png":            "./img.png",
		"./my%20image.png":     "./my%20image.png",
		"./my image.png":       "<./my image.png>",
		"./image (1).png":      "<./image (1).png>",
		"./image(1).png":       "./image(1).png",
		"./image(.png":         "<./image(.png>",
		"./image).png":         "<./image).png>",
		`./image\).png`:        `./image\).png`,
		"./a<b>.png":           `<./a\<b\>.png>`,
		"https://example.com/": "https://example.com/",
	}

	for destination, expected := range testCases {
		if g := string(formatDestination([]byte(destination))); g != expected {
			t.Errorf("formatDestination(%q): expected %q, got %q", destination, expected, g)
		}
	}
}
