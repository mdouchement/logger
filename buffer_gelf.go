package logger

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// A BufferGELF is a buffer used to build GELF payload.
type BufferGELF struct {
	buf *bytes.Buffer
}

// NewBufferGELF returns a new BufferGELF.
func NewBufferGELF() *BufferGELF {
	buf := bytes.NewBuffer(make([]byte, 0, 256))
	buf.WriteString(`{"version":"1.1"`)
	return &BufferGELF{
		buf: buf,
	}
}

// Host adds the host to the GELF buffer.
func (b *BufferGELF) Host(h string) {
	b.key("host")
	b.string(h, true)
}

// Level adds the level to the GELF buffer.
func (b *BufferGELF) Level(l int32) {
	b.key("level")
	b.buf.WriteString(fmt.Sprint(l))
}

// Message adds the short_message/full_message to the GELF buffer.
func (b *BufferGELF) Message(m string) {
	if i := strings.IndexRune(m, '\n'); i > 0 {
		// If there are newlines in the message, use the first line
		// for the short_message and set the full_message to the
		// original input. If the input has no newlines, stick the
		// whole thing in short_message.
		b.key("short_message")
		b.string(m[:i], true)
		b.key("full_message")
		b.string(m, true)
		return
	}

	b.key("short_message")
	b.string(m, true)
}

// Timestamp adds the timestamp to the GELF buffer.
func (b *BufferGELF) Timestamp(t time.Time) {
	b.key("timestamp")
	b.buf.WriteString(strconv.FormatFloat(
		float64(t.UnixNano())/1e9, // Unix epoch timestamp in seconds
		'f',
		-1,
		64,
	))
}

// Add adds any key/value to the GELF buffer.
func (b *BufferGELF) Add(k string, v any) {
	// skip id
	if k == "id" || k == "_id" {
		return
	}

	b.key("_" + k)

	// otherwise convert if necessary
	switch value := v.(type) {
	case time.Time:
		b.string(value.Format(time.RFC3339), true)
	case uint, uint8, uint16, uint32,
		int, int8, int16, int32, int64:
		b.buf.WriteString(fmt.Sprint(value))
	case uint64:
		// NOTE: uint64 is not supported by graylog due to java limitation
		//       so we're sending them as double for the time being
		b.buf.WriteString(strconv.FormatFloat(float64(value), 'f', -1, 64))
	case float32:
		b.buf.WriteString(strconv.FormatFloat(float64(value), 'f', -1, 32))
	case float64:
		b.buf.WriteString(strconv.FormatFloat(value, 'f', -1, 64))
	case bool:
		b.buf.WriteString(fmt.Sprintf(`"%t"`, value))
	case string:
		b.string(value, true)
	default:
		b.string(fmt.Sprint(value), true)
	}
}

// Complete returns the completed GELF payload with a `\n' when ln is true.
func (b *BufferGELF) Complete(ln bool) []byte {
	b.buf.WriteString("}")
	if ln {
		b.buf.WriteString("\n")
	}

	return b.buf.Bytes()
}

// Bytes returns the bytes of the current GELF payload state.
func (b *BufferGELF) Bytes() []byte {
	return b.buf.Bytes()
}

func (b *BufferGELF) key(k string) {
	b.buf.WriteString(",")
	b.string(k, true)
	b.buf.WriteString(":")
}

// based on https://cs.opensource.google/go/go/+/refs/tags/go1.22.0:src/encoding/json/encode.go;l=956
func (b *BufferGELF) string(src string, escapeHTML bool) {
	buf := b.buf
	buf.WriteByte('"')
	start := 0
	for i := 0; i < len(src); {
		if b := src[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] || (!escapeHTML && safeSet[b]) {
				i++
				continue
			}
			buf.WriteString(src[start:i])
			switch b {
			case '\\', '"':
				buf.WriteByte('\\')
				buf.WriteByte(b)
			case '\b':
				buf.WriteByte('\\')
				buf.WriteByte('b')
			case '\f':
				buf.WriteByte('\\')
				buf.WriteByte('f')
			case '\n':
				buf.WriteByte('\\')
				buf.WriteByte('n')
			case '\r':
				buf.WriteByte('\\')
				buf.WriteByte('r')
			case '\t':
				buf.WriteByte('\\')
				buf.WriteByte('t')
			default:
				// This encodes bytes < 0x20 except for \b, \f, \n, \r and \t.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				buf.WriteByte('\\')
				buf.WriteByte('u')
				buf.WriteByte('0')
				buf.WriteByte('0')
				buf.WriteByte(hex[b>>4])
				buf.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		// TODO(https://go.dev/issue/56948): Use generic utf8 functionality.
		// For now, cast only a small portion of byte slices to a string
		// so that it can be stack allocated. This slows down []byte slightly
		// due to the extra copy, but keeps string performance roughly the same.
		n := len(src) - i
		if n > utf8.UTFMax {
			n = utf8.UTFMax
		}
		c, size := utf8.DecodeRuneInString(string(src[i : i+n]))
		if c == utf8.RuneError && size == 1 {
			b.buf.WriteString(src[start:i])
			b.buf.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See https://en.wikipedia.org/wiki/JSON#Safety.
		if c == '\u2028' || c == '\u2029' {
			b.buf.WriteString(src[start:i])
			b.buf.WriteByte('\\')
			b.buf.WriteByte('u')
			b.buf.WriteByte('2')
			b.buf.WriteByte('0')
			b.buf.WriteByte('2')
			b.buf.WriteByte(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	b.buf.WriteString(src[start:])
	b.buf.WriteByte('"')
}

const hex = "0123456789abcdef"

// safeSet holds the value true if the ASCII character with the given array
// position can be represented inside a JSON string without any further
// escaping.
//
// All values are true except for the ASCII control characters (0-31), the
// double quote ("), and the backslash character ("\").
var safeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      true,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      true,
	'=':      true,
	'>':      true,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}

// htmlSafeSet holds the value true if the ASCII character with the given
// array position can be safely represented inside a JSON string, embedded
// inside of HTML <script> tags, without any additional escaping.
//
// All values are true except for the ASCII control characters (0-31), the
// double quote ("), the backslash character ("\"), HTML opening and closing
// tags ("<" and ">"), and the ampersand ("&").
var htmlSafeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      false,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      false,
	'=':      true,
	'>':      false,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}
