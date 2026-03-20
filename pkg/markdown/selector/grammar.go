package selector

import (
	"fmt"
	"strings"
	"unicode"
)

type parser struct {
	input string
	pos   int
}

func newParser(input string) *parser {
	return &parser{input: strings.TrimSpace(input)}
}

func (p *parser) parseSelectorList() (SelectorList, error) {
	if p.pos >= len(p.input) {
		return nil, fmt.Errorf("empty selector")
	}

	var result SelectorList

	sel, err := p.parseSelector()
	if err != nil {
		return nil, err
	}
	result = append(result, sel)

	for {
		p.skipSpaces()
		if p.pos >= len(p.input) {
			break
		}
		if p.input[p.pos] != ',' {
			break
		}
		p.pos++ // consume ','
		p.skipSpaces()

		sel, err = p.parseSelector()
		if err != nil {
			return nil, err
		}
		result = append(result, sel)
	}

	return result, nil
}

func (p *parser) parseSelector() (*compoundSelector, error) {
	simple, err := p.parseSimpleSelector()
	if err != nil {
		return nil, err
	}

	cs := &compoundSelector{
		left:  nil,
		right: simple,
	}

	for {
		savedPos := p.pos

		hasSpace := p.skipSpaces()

		if p.pos >= len(p.input) || p.input[p.pos] == ',' {
			break
		}

		var combinator combinatorKind
		if p.pos < len(p.input) && p.input[p.pos] == '>' {
			p.pos++ // consume '>'
			p.skipSpaces()
			combinator = combinatorChild
		} else if p.pos < len(p.input) && p.input[p.pos] == '~' {
			p.pos++ // consume '~'
			p.skipSpaces()
			combinator = combinatorSibling
		} else if hasSpace {
			combinator = combinatorDescendant
		} else {
			p.pos = savedPos
			break
		}

		simple, err = p.parseSimpleSelector()
		if err != nil {
			return nil, err
		}

		cs = &compoundSelector{
			left:       cs,
			combinator: combinator,
			right:      simple,
		}
	}

	return cs, nil
}

func (p *parser) parseSimpleSelector() (simpleSelector, error) {
	var sel simpleSelector

	// Type selector: must start with alpha
	if p.pos < len(p.input) && isAlpha(p.input[p.pos]) {
		start := p.pos
		for p.pos < len(p.input) && isIdentChar(p.input[p.pos]) {
			p.pos++
		}
		typeStr := p.input[start:p.pos]
		switch typeStr {
		case "h1", "h2", "h3", "h4", "h5", "h6", "p", "blockquote", "code", "ul", "ol", "table":
			sel.typeStr = typeStr
		default:
			p.pos = start // backtrack: not a known type selector
		}
	}

	// ID selector: #ident
	if p.pos < len(p.input) && p.input[p.pos] == '#' {
		p.pos++ // consume '#'
		ident, err := p.parseIdent()
		if err != nil {
			return sel, fmt.Errorf("expected identifier after '#': %w", err)
		}
		sel.id = ident
	}

	// Attribute selector: [name="value"]
	if p.pos < len(p.input) && p.input[p.pos] == '[' {
		p.pos++ // consume '['
		name, err := p.parseIdent()
		if err != nil {
			return sel, fmt.Errorf("expected attribute name in '[...]': %w", err)
		}
		if p.pos >= len(p.input) || p.input[p.pos] != '=' {
			return sel, fmt.Errorf("expected '=' in attribute selector at pos %d", p.pos)
		}
		p.pos++ // consume '='
		val, err := p.parseQuotedString()
		if err != nil {
			return sel, fmt.Errorf("expected quoted string in attribute selector: %w", err)
		}
		if p.pos >= len(p.input) || p.input[p.pos] != ']' {
			return sel, fmt.Errorf("expected ']' to close attribute selector at pos %d", p.pos)
		}
		p.pos++ // consume ']'
		sel.attrName = name
		sel.attrVal = val
	}

	// :contains("pattern")
	if p.pos < len(p.input) && p.input[p.pos] == ':' {
		savedPos := p.pos
		p.pos++ // consume ':'
		ident, err := p.parseIdent()
		if err != nil || ident != "contains" {
			p.pos = savedPos // backtrack
		} else {
			if p.pos >= len(p.input) || p.input[p.pos] != '(' {
				return sel, fmt.Errorf("expected '(' after ':contains' at pos %d", p.pos)
			}
			p.pos++ // consume '('
			val, err := p.parseQuotedString()
			if err != nil {
				return sel, fmt.Errorf("expected quoted string in :contains(): %w", err)
			}
			if p.pos >= len(p.input) || p.input[p.pos] != ')' {
				return sel, fmt.Errorf("expected ')' to close :contains() at pos %d", p.pos)
			}
			p.pos++ // consume ')'
			sel.contains = val
		}
	}

	if sel.typeStr == "" && sel.id == "" && sel.attrName == "" && sel.contains == "" {
		if p.pos >= len(p.input) {
			return sel, fmt.Errorf("expected selector but got EOF")
		}
		return sel, fmt.Errorf("expected selector at pos %d (found %q)", p.pos, p.input[p.pos:])
	}

	return sel, nil
}

func (p *parser) parseIdent() (string, error) {
	start := p.pos
	for p.pos < len(p.input) && isIdentChar(p.input[p.pos]) {
		p.pos++
	}
	if p.pos == start {
		if p.pos >= len(p.input) {
			return "", fmt.Errorf("expected identifier but got EOF")
		}
		return "", fmt.Errorf("expected identifier but got %q at pos %d", string(p.input[p.pos]), p.pos)
	}
	return p.input[start:p.pos], nil
}

func (p *parser) parseQuotedString() (string, error) {
	if p.pos >= len(p.input) {
		return "", fmt.Errorf("expected quoted string but got EOF")
	}
	quote := p.input[p.pos]
	if quote != '"' && quote != '\'' {
		return "", fmt.Errorf("expected quote character but got %q at pos %d", string(quote), p.pos)
	}
	p.pos++ // consume opening quote

	var buf strings.Builder
	for p.pos < len(p.input) {
		c := p.input[p.pos]
		if c == '\\' && p.pos+1 < len(p.input) {
			p.pos++
			buf.WriteByte(p.input[p.pos])
			p.pos++
			continue
		}
		if c == quote {
			p.pos++ // consume closing quote
			return buf.String(), nil
		}
		buf.WriteByte(c)
		p.pos++
	}
	return "", fmt.Errorf("unterminated string starting at pos %d", p.pos)
}

// skipSpaces skips whitespace and returns true if any was consumed.
func (p *parser) skipSpaces() bool {
	start := p.pos
	for p.pos < len(p.input) && unicode.IsSpace(rune(p.input[p.pos])) {
		p.pos++
	}
	return p.pos > start
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isIdentChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_'
}
