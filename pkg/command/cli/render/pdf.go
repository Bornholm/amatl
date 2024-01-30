package render

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/yuin/goldmark/text"
)

func PDF() *cli.Command {
	return &cli.Command{
		Name: "pdf",
		Flags: withCommonFlags(&cli.StringFlag{
			Name:     "output",
			Required: true,
			Value:    "",
		},
		),
		Action: func(ctx *cli.Context) error {
			for _, filename := range ctx.Args().Slice() {
				basePath, err := filepath.Abs(filepath.Dir(filename))
				if err != nil {
					return errors.WithStack(err)
				}

				source, err := os.ReadFile(filename)
				if err != nil {
					return errors.WithStack(err)
				}

				reader := text.NewReader(source)

				parse := newParser(basePath, false)
				render := newMarkdownRenderer()

				document := parse.Parse(reader)

				var buf bytes.Buffer
				if err := render.Render(&buf, source, document); err != nil {
					return errors.WithStack(err)
				}

				source = buf.Bytes()
				reader = text.NewReader(source)

				withToc := isTocEnabled(ctx)

				parse = newParser(basePath, withToc)
				document = parse.Parse(reader)

				render = newHTMLRenderer()

				var html bytes.Buffer
				var pdf []byte
				if err := render.Render(&html, source, document); err != nil {
					return errors.WithStack(err)
				}

				c, _ := chromedp.NewContext(context.Background())
				if err := chromedp.Run(c, printToPDF(html.Bytes(), &pdf)); err != nil {
					fmt.Println("oops")
					return errors.WithStack(err)
				}

				if err := os.WriteFile(ctx.String("output"), pdf, 0o644); err != nil {
					return errors.WithStack(err)
				}
				fmt.Println("wrote ", ctx.String("output"))

			}

			return nil
		},
	}
}

func printToPDF(html []byte, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			lctx, cancel := context.WithCancel(ctx)
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(1)
			chromedp.ListenTarget(lctx, func(ev interface{}) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					// It's a good habit to remove the event listener if we don't need it anymore.
					cancel()
					wg.Done()
				}
			})
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			if err := page.SetDocumentContent(frameTree.Frame.ID, string(html)).Do(ctx); err != nil {
				return err
			}
			wg.Wait()
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
