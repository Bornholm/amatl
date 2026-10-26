package container

import (
	"bytes"
	"testing"

	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown/node"
	"github.com/andreyvit/diff"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

type testCase struct {
	Name     string
	Source   string
	Expected string
}

var testCases = []testCase{
	{
		Name: "simple",
		Source: `:::card
Hello **world**
:::
`,
		Expected: `<div class="amatl-card">
<p>Hello <strong>world</strong></p>
</div>
`,
	},
	{
		Name: "attributes",
		Source: `:::card{#intro .highlight .large count=3 style="color: red" role=note aria-label="Intro" data-foo=bar}
Hello
:::
`,
		Expected: `<div class="amatl-card highlight large" id="intro" data-count="3" style="color: red" role="note" aria-label="Intro" data-foo="bar">
<p>Hello</p>
</div>
`,
	},
	{
		Name: "escaped attribute value",
		Source: `:::card{title="Say \"hi\" \\o/"}
Hello
:::
`,
		Expected: `<div class="amatl-card" title="Say &quot;hi&quot; \o/">
<p>Hello</p>
</div>
`,
	},
	{
		Name: "semantic tag",
		Source: `:::section{.hero}
## Title
:::
`,
		Expected: `<section class="amatl-section hero">
<h2>Title</h2>
</section>
`,
	},
	{
		Name: "nested with same fence length",
		Source: `:::columns
:::card
One
:::
:::card
Two
:::
:::

After
`,
		Expected: `<div class="amatl-columns">
<div class="amatl-card">
<p>One</p>
</div>
<div class="amatl-card">
<p>Two</p>
</div>
</div>
<p>After</p>
`,
	},
	{
		Name: "nested with longer outer fence",
		Source: `::::columns
:::card
One
:::
::::
`,
		Expected: `<div class="amatl-columns">
<div class="amatl-card">
<p>One</p>
</div>
</div>
`,
	},
	{
		Name: "longer closing fence closes the innermost container",
		Source: `::::outer
:::inner
One
::::
Two
::::
`,
		Expected: `<div class="amatl-outer">
<div class="amatl-inner">
<p>One</p>
</div>
<p>Two</p>
</div>
`,
	},
	{
		Name: "indented nested containers",
		Source: `:::columns
  :::card
  One
  :::
:::
`,
		Expected: `<div class="amatl-columns">
<div class="amatl-card">
<p>One</p>
</div>
</div>
`,
	},
	{
		Name: "fence in list item closes the inner container",
		Source: `:::outer
- item

  :::inner
  One
  :::
:::
After
`,
		Expected: `<div class="amatl-outer">
<ul>
<li>
<p>item</p>
<div class="amatl-inner">
<p>One</p>
</div>
</li>
</ul>
</div>
<p>After</p>
`,
	},
	{
		Name: "dedented fence closes the outer container",
		Source: `:::outer
- item

  :::inner
  One

:::
After
`,
		Expected: `<div class="amatl-outer">
<ul>
<li>
<p>item</p>
<div class="amatl-inner">
<p>One</p>
</div>
</li>
</ul>
</div>
<p>After</p>
`,
	},
	{
		Name: "dedented fence after lazy continuation closes the outer container",
		Source: `:::outer
- item

  :::inner
  One
:::
After
`,
		Expected: `<div class="amatl-outer">
<ul>
<li>
<p>item</p>
<div class="amatl-inner">
<p>One</p>
</div>
</li>
</ul>
</div>
<p>After</p>
`,
	},
	{
		Name: "dedented fence closes the container matching its indentation",
		Source: `:::outer
- a

  :::mid
  - b

    :::inner
    One
  :::
  Two
:::
After
`,
		Expected: `<div class="amatl-outer">
<ul>
<li>
<p>a</p>
<div class="amatl-mid">
<ul>
<li>
<p>b</p>
<div class="amatl-inner">
<p>One</p>
</div>
</li>
</ul>
</div>
<p>Two</p>
</li>
</ul>
</div>
<p>After</p>
`,
	},
	{
		Name: "inside blockquote",
		Source: `:::outer
> :::inner
> One
> :::
:::
`,
		Expected: `<div class="amatl-outer">
<blockquote>
<div class="amatl-inner">
<p>One</p>
</div>
</blockquote>
</div>
`,
	},
	{
		Name:   "tab indented fences in list item",
		Source: "- item\n\n\t:::card\n\tOne\n\t:::\n\nAfter\n",
		Expected: `<ul>
<li>
<p>item</p>
<div class="amatl-card">
<p>One</p>
</div>
</li>
</ul>
<p>After</p>
`,
	},
	{
		Name:   "CRLF line endings",
		Source: ":::card\r\nOne\r\n:::\r\nAfter\r\n",
		Expected: `<div class="amatl-card">
<p>One</p>
</div>
<p>After</p>
`,
	},
	{
		Name: "interrupts paragraph",
		Source: `Before
:::card
Inside
:::
After
`,
		Expected: `<p>Before</p>
<div class="amatl-card">
<p>Inside</p>
</div>
<p>After</p>
`,
	},
	{
		Name: "unclosed container ends with document",
		Source: `:::card
Inside
`,
		Expected: `<div class="amatl-card">
<p>Inside</p>
</div>
`,
	},
	{
		Name: "inside list item",
		Source: `- item

  :::card
  Inside
  :::
`,
		Expected: `<ul>
<li>
<p>item</p>
<div class="amatl-card">
<p>Inside</p>
</div>
</li>
</ul>
`,
	},
	{
		Name:   "fenced code is not parsed",
		Source: "```\n:::card\n```\n",
		Expected: `<pre><code>:::card
</code></pre>
`,
	},
	{
		Name: "not a container",
		Source: `:: card

:::

:::1card

:::card{unclosed
`,
		Expected: `<p>:: card</p>
<p>:::</p>
<p>:::1card</p>
<p>:::card{unclosed</p>
`,
	},
}

func TestHTMLRenderer(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			html := renderHTML(t, []byte(tc.Source))

			if e, g := tc.Expected, html; e != g {
				t.Errorf("unexpected html:\n%s", diff.LineDiff(e, g))
			}
		})
	}
}

func TestMarkdownRendererRoundTrip(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			source := []byte(tc.Source)

			md := goldmark.New(goldmark.WithExtensions(extension.GFM, &Extension{}))
			doc := md.Parser().Parse(text.NewReader(source))

			renderer := markdown.NewRenderer()
			renderer.AddMarkdownOptions(
				markdown.WithNodeRenderers(node.Renderers()),
				markdown.WithNodeRenderer(Kind, node.WithLineSpacingBefore(&MarkdownRenderer{}, 2)),
			)

			var rendered bytes.Buffer
			if err := renderer.Render(&rendered, source, doc); err != nil {
				t.Fatalf("%+v", err)
			}

			if e, g := renderHTML(t, source), renderHTML(t, rendered.Bytes()); e != g {
				t.Errorf("html differs after markdown round trip:\n%s\n\nmarkdown:\n%s", diff.LineDiff(e, g), rendered.String())
			}
		})
	}
}

func TestMarshalAttributes(t *testing.T) {
	source := []byte(`:::card{#intro .a .b title="Say \"hi\"" count=3}
:::
`)

	md := goldmark.New(goldmark.WithExtensions(&Extension{}))
	doc := md.Parser().Parse(text.NewReader(source))

	container, ok := doc.FirstChild().(*Node)
	if !ok {
		t.Fatalf("expected *container.Node, got %T", doc.FirstChild())
	}

	if e, g := `#intro .a .b title="Say \"hi\"" count="3"`, marshalAttributes(container.Attributes()); e != g {
		t.Errorf("marshalAttributes: expected '%s', got '%s'", e, g)
	}
}

func TestSourceRange(t *testing.T) {
	source := []byte("Intro\n\n:::columns\n:::card\nText\n:::\n:::\n\n:::unclosed\nEnd\n")

	md := goldmark.New(goldmark.WithExtensions(&Extension{}))
	doc := md.Parser().Parse(text.NewReader(source))

	columns, ok := doc.FirstChild().NextSibling().(*Node)
	if !ok {
		t.Fatalf("expected *container.Node, got %T", doc.FirstChild().NextSibling())
	}

	start, stop := columns.SourceRange()
	if e, g := ":::columns\n:::card\nText\n:::\n:::\n", string(source[start:stop]); e != g {
		t.Errorf("SourceRange: expected %q, got %q", e, g)
	}

	card, ok := columns.FirstChild().(*Node)
	if !ok {
		t.Fatalf("expected *container.Node, got %T", columns.FirstChild())
	}

	start, stop = card.SourceRange()
	if e, g := ":::card\nText\n:::\n", string(source[start:stop]); e != g {
		t.Errorf("SourceRange: expected %q, got %q", e, g)
	}

	unclosed, ok := columns.NextSibling().(*Node)
	if !ok {
		t.Fatalf("expected *container.Node, got %T", columns.NextSibling())
	}

	start, stop = unclosed.SourceRange()
	if e, g := ":::unclosed\nEnd\n", string(source[start:]); e != g || stop != -1 {
		t.Errorf("SourceRange: expected %q and -1, got %q and %d", e, g, stop)
	}
}

func renderHTML(t *testing.T, source []byte) string {
	t.Helper()

	md := goldmark.New(goldmark.WithExtensions(extension.GFM, &Extension{}))

	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		t.Fatalf("%+v", err)
	}

	return buf.String()
}
