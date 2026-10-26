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

func renderHTML(t *testing.T, source []byte) string {
	t.Helper()

	md := goldmark.New(goldmark.WithExtensions(extension.GFM, &Extension{}))

	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		t.Fatalf("%+v", err)
	}

	return buf.String()
}
