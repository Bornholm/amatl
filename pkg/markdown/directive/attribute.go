package directive

import (
	"fmt"
	"strconv"
	"strings"
)

// NormalizeAttributeValue converts an attribute value as parsed by goldmark
// (which can be a []byte, a float64, a bool, nil or a nested structure) to
// its textual representation
func NormalizeAttributeValue(value any) []byte {
	return []byte(AttributeValueString(value))
}

// AttributeValueString returns the textual representation of an attribute value
func AttributeValueString(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case []byte:
		return string(typed)
	case string:
		return typed
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(typed)
	default:
		return fmt.Sprintf("%v", typed)
	}
}

// QuoteAttributeValue quotes the given value so it can be parsed back
// as an attribute string value
func QuoteAttributeValue(value string) string {
	var sb strings.Builder

	sb.WriteRune('"')
	for _, r := range value {
		switch r {
		case '"', '\\':
			sb.WriteRune('\\')
			sb.WriteRune(r)
		case '\n':
			sb.WriteString(`\n`)
		case '\r':
			sb.WriteString(`\r`)
		case '\t':
			sb.WriteString(`\t`)
		default:
			sb.WriteRune(r)
		}
	}
	sb.WriteRune('"')

	return sb.String()
}
