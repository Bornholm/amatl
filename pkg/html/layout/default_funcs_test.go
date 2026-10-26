package layout

import (
	"html/template"
	"testing"
)

func TestHTMLTextContent(t *testing.T) {
	body := template.HTML(`<h1 id="title">Hello <em>world</em></h1><p>After</p><h1>Second</h1>`)

	text, err := htmlTextContent(body, "h1:first-of-type")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if e, g := "Hello world", text; e != g {
		t.Errorf("htmlTextContent: expected '%s', got '%s'", e, g)
	}

	text, err = htmlTextContent(body, "h1")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if e, g := "Hello worldSecond", text; e != g {
		t.Errorf("htmlTextContent: expected '%s', got '%s'", e, g)
	}
}

func TestHTMLRemove(t *testing.T) {
	body := template.HTML(`<style>p{}</style><h1>Title</h1><p>Content</p>`)

	html, err := htmlRemove(body, "h1")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if e, g := template.HTML(`<style>p{}</style><p>Content</p>`), html; e != g {
		t.Errorf("htmlRemove: expected '%s', got '%s'", e, g)
	}
}

func TestHTMLRemoveKeepsTopLevelNodes(t *testing.T) {
	body := template.HTML(`<!-- lead --><h1>Title</h1><p>Content</p><!-- tail -->`)

	html, err := htmlRemove(body, "h1")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if e, g := template.HTML(`<!-- lead --><p>Content</p><!-- tail -->`), html; e != g {
		t.Errorf("htmlRemove: expected '%s', got '%s'", e, g)
	}
}

func TestHTMLAddAttr(t *testing.T) {
	body := template.HTML(`<p class="a">One</p><p>Two</p>`)

	html, err := htmlAddAttr(body, "p", "class", "b")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if e, g := template.HTML(`<p class="a b">One</p><p class="b">Two</p>`), html; e != g {
		t.Errorf("htmlAddAttr: expected '%s', got '%s'", e, g)
	}
}
