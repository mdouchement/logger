package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/mdouchement/logger/syslog"
)

type (
	SlogGELFOption struct {
		Level    slog.Level
		Hostname string
	}

	// A SlogGELFHandler is GELF formatter for log/slog.
	SlogGELFHandler struct {
		writer io.Writer

		level    slog.Level
		hostname string

		prefix string
		group  string
		attrs  []slog.Attr
	}
)

// NewSlogGELFHandler returns a new SlogGELFHandler.
func NewSlogGELFHandler(w io.Writer, o *SlogGELFOption) *SlogGELFHandler {
	var err error

	if o == nil {
		o = &SlogGELFOption{Level: slog.LevelInfo}
	}
	if o.Hostname == "" {
		o.Hostname, err = os.Hostname()
		if err != nil {
			o.Hostname = "localhost"
		}
	}

	return &SlogGELFHandler{
		level:    o.Level,
		hostname: o.Hostname,
		writer:   w,
	}
}

// Enabled reports whether the handler handles records at the given level.
func (f *SlogGELFHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= f.level
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
func (h *SlogGELFHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	var prefix string
	if h.group != "" {
		prefix = h.group + "."
	}

	nh := h.Clone() // Since h is cloned, it's just read, never edited so it's safe.
	for _, attr := range attrs {
		if attr.Key == KeyPrefix {
			nh.prefix += attr.Value.String()
			continue
		}

		attr.Key = prefix + attr.Key
		nh.attrs = append(nh.attrs, attr)
	}

	return nh
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
func (h *SlogGELFHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	nh := h.Clone() // Since h is cloned, it's just read, never edited so it's safe.
	if h.group != "" {
		nh.group = h.group + "." + name
	} else {
		nh.group = name
	}

	return nh
}

// Handle handles the Record.
func (h *SlogGELFHandler) Handle(_ context.Context, record slog.Record) error {
	m := make(map[string]any, len(h.attrs)+record.NumAttrs()+10)

	// Process handler's groups/attrs.
	for _, attr := range h.attrs {
		gelf(m, attr)
	}

	// Process record's groups/attrs.
	if record.NumAttrs() > 0 {
		var prefix string
		if h.group != "" {
			prefix = h.group + "."
		}

		record.Attrs(func(attr slog.Attr) bool {
			gelfrecord(m, prefix, attr)
			return true
		})
	}

	// Main fields.
	m["version"] = "1.1"
	m["host"] = h.hostname
	m["timestamp"] = float64(record.Time.UnixNano()) / 1e9 // Unix epoch timestamp in seconds with decimals for nanoseconds.
	m["level"] = h.priorities(record.Level)

	if h.prefix != "" {
		record.Message = fmt.Sprintf("%s %s", h.prefix, record.Message)
	}
	m["short_message"] = record.Message
	// If there are newlines in the message, use the first line
	// for the short_message and set the full_message to the
	// original input. If the input has no newlines, stick the
	// whole thing in short_message.
	if i := strings.IndexRune(record.Message, '\n'); i > 0 {
		m["short_message"] = record.Message[:i]
		m["full_message"] = record.Message
	}

	m["_level_name"] = record.Level.String()
	// m["_file"]
	// m["_line"]

	payload, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = h.writer.Write(append(payload, '\n'))
	return err
}

func (h *SlogGELFHandler) Clone() *SlogGELFHandler {
	nh := &SlogGELFHandler{
		level:    h.level,
		hostname: h.hostname,
		writer:   h.writer,
		prefix:   h.prefix,
		group:    h.group,
		attrs:    make([]slog.Attr, len(h.attrs)),
	}
	copy(nh.attrs, h.attrs)

	return nh
}

func (SlogGELFHandler) priorities(level slog.Level) int32 {
	var p syslog.Priority

	switch level {
	case slog.LevelError:
		p = syslog.LOG_ERR
	case slog.LevelWarn:
		p = syslog.LOG_WARNING
	case slog.LevelInfo:
		p = syslog.LOG_INFO
	case slog.LevelDebug:
		p = syslog.LOG_DEBUG
	}

	return int32(p)
}

func gelfrecord(m map[string]any, p string, attr slog.Attr) {
	if attr.Value.Kind() != slog.KindGroup {
		attr.Key = p + attr.Key
		gelf(m, attr)
		return
	}

	p += attr.Key + "."
	for _, a := range attr.Value.Group() {
		if a.Value.Kind() == slog.KindGroup {
			gelfrecord(m, p, a)
			continue
		}
		a.Key = p + a.Key
		gelf(m, a)
	}
}

func gelf(m map[string]any, attr slog.Attr) {
	// skip id if present
	if attr.Key == "id" || attr.Key == "_id" {
		return
	}

	// otherwise convert if necessary
	switch value := attr.Value.Any().(type) {
	case time.Time:
		m["_"+attr.Key] = value.Format(time.RFC3339)
	case uint, uint8, uint16, uint32,
		int, int8, int16, int32, int64:
		m["_"+attr.Key] = value
	case uint64, float32, float64:
		// NOTE: uint64 is not supported by graylog due to java limitation
		//       so we're sending them as double for the time being
		m["_"+attr.Key] = value
	case bool:
		m["_"+attr.Key] = fmt.Sprintf("%t", value)
	case string:
		m["_"+attr.Key] = value
	default:
		m["_"+attr.Key] = fmt.Sprint(value)
	}
}
