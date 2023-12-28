package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/mdouchement/logger/syslog"
)

type (
	// A SlogGELFOption holds SlogGELFHandler's options.
	SlogGELFOption struct {
		Level    slog.Level
		Hostname string
	}

	// A SlogGELFHandler is GELF formatter for log/slog.
	SlogGELFHandler struct {
		opt    *SlogGELFOption
		writer io.Writer

		parent *SlogGELFHandler
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
		opt:    o,
		writer: w,
	}
}

// Enabled reports whether the handler handles records at the given level.
func (h *SlogGELFHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.opt.Level
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
func (h *SlogGELFHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	nh := h.Clone() // Since h is cloned, it's just read, never edited so it's safe.
	nh.attrs = make([]slog.Attr, 0, len(attrs))
	for _, attr := range attrs {
		if attr.Key == KeyPrefix {
			nh.prefix += attr.Value.String()
			continue
		}

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
	nh.group = name
	return nh
}

// Handle handles the Record.
func (h *SlogGELFHandler) Handle(_ context.Context, record slog.Record) error {
	gelf := NewBufferGELF()

	// Get all parents in a list.
	ilineage := make([]*SlogGELFHandler, 0, 100)
	p := h
	for p != nil {
		ilineage = append(ilineage, p)
		p = p.parent
	}

	// Process handler's groups/attrs.
	var gprefix string
	for i := len(ilineage) - 1; i >= 0; i-- {
		p = ilineage[i]
		if p.group != "" {
			gprefix += p.group + delimiter
		}

		for _, attr := range ilineage[i].attrs {
			gelf.Add(gprefix+attr.Key, attr.Value.Any())
		}
	}

	// Process record's groups/attrs.
	if record.NumAttrs() > 0 {
		record.Attrs(func(attr slog.Attr) bool {
			gelfrecord(gelf, gprefix, attr)
			return true
		})
	}

	// Main fields.
	gelf.Host(h.opt.Hostname)
	gelf.Timestamp(record.Time)
	gelf.Level(h.priorities(record.Level))

	if h.prefix != "" {
		record.Message = fmt.Sprintf("%s %s", h.prefix, record.Message)
	}
	gelf.Message(record.Message)

	gelf.Add("level_name", record.Level.String())
	// m["_file"]
	// m["_line"]

	_, err := h.writer.Write(gelf.Complete(true))
	return err
}

// Clone clones the entry, it creates a new instance and linking the parent to it.
func (h *SlogGELFHandler) Clone() *SlogGELFHandler {
	nh := &SlogGELFHandler{
		parent: h,
		opt:    h.opt,
		prefix: h.prefix,
		writer: h.writer,
	}

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

func gelfrecord(gelf *BufferGELF, p string, attr slog.Attr) {
	if attr.Value.Kind() != slog.KindGroup {
		gelf.Add(p+attr.Key, attr.Value.Any())
		return
	}

	p += attr.Key + delimiter
	for _, a := range attr.Value.Group() {
		if a.Value.Kind() == slog.KindGroup {
			gelfrecord(gelf, p, a)
			continue
		}
		gelf.Add(p+a.Key, a.Value.Any())
	}
}
