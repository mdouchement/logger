package logger

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const buflen = 256

// A BufferGELF is a buffer used to build GELF payload.
type BufferGELF struct {
	buf *bytes.Buffer
}

// NewBufferGELF returns a new BufferGELF.
func NewBufferGELF() *BufferGELF {
	return NewBufferGELFs(buflen)
}

// NewBufferGELFs returns a new BufferGELF with buffer of size.
func NewBufferGELFs(size int) *BufferGELF {
	buf := bytes.NewBuffer(make([]byte, 0, size))
	buf.WriteString(`{"version":"1.1"`)
	return &BufferGELF{
		buf: buf,
	}
}

// Host adds the host to the GELF buffer.
func (b *BufferGELF) Host(h string) {
	b.key("host")
	b.buf.WriteString(strconv.QuoteToGraphic(h))
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
		b.buf.WriteString(strconv.QuoteToGraphic(m[:i]))
		b.key("full_message")
		b.buf.WriteString(strconv.QuoteToGraphic(m))
		return
	}

	b.key("short_message")
	b.buf.WriteString(strconv.QuoteToGraphic(m))
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
		b.buf.WriteString(strconv.QuoteToGraphic(value.Format(time.RFC3339)))
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
		b.buf.WriteString(strconv.QuoteToGraphic(value))
	default:
		b.buf.WriteString(strconv.QuoteToGraphic(fmt.Sprint(value)))
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

// Clone clones the current state of the BufferGELF.
func (b *BufferGELF) Clone() *BufferGELF {
	nb := NewBufferGELFs(b.buf.Cap())
	nb.buf.Reset() // Remove initialization data.
	nb.buf.Write(b.Bytes())
	return nb
}

func (b *BufferGELF) key(k string) {
	b.buf.WriteString(",")
	b.buf.WriteString(strconv.QuoteToGraphic(k))
	b.buf.WriteString(":")
}
