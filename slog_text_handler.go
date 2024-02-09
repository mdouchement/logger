package logger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"sort"
	"strings"
)

type (
	// SlogTextOption holds SlogTextHandler's options.
	SlogTextOption struct {
		// Set the logger's level.
		Level slog.Level

		// Set to true to bypass checking for a TTY before outputting colors.
		ForceColors bool

		// Force disabling colors. For a TTY colors are enabled by default.
		DisableColors bool

		// Force formatted layout, even for non-TTY output.
		ForceFormatting bool

		// Disable timestamp logging. useful when output is redirected to logging
		// system that already adds timestamps.
		DisableTimestamp bool

		// Disable the conversion of the log levels to uppercase
		DisableUppercase bool

		// Enable logging the full timestamp when a TTY is attached instead of just
		// the time passed since beginning of execution.
		FullTimestamp bool

		// Timestamp format to use for display when a full timestamp is printed.
		TimestampFormat string

		// The fields are sorted by default for a consistent output. For applications
		// that log extremely frequently and don't use the JSON formatter this may not
		// be desired.
		DisableSorting bool

		// Wrap empty fields in quotes if true.
		QuoteEmptyFields bool

		// Can be set to the override the default quoting character "
		// with something else. For example: ', or `.
		QuoteCharacter string

		// Pad msg field with spaces on the right for display.
		// The value for this parameter will be the size of padding.
		// Its default value is zero, which means no padding will be applied for msg.
		SpacePadding int

		// Regexp to find prefix to be colored.
		// If not defined, no coloration is applied.
		// e.g. Use `^(\[.*?\])\s` to colorize prefix for message like "[prefix#1][prefix#2] The message"
		PrefixRE *regexp.Regexp

		// ValueFormatter is the format of the value when logs are pretty printed.
		// The default value is `%v'. You can use `%+v' to print the stacktrace of github.com/pkg/errors.
		ValueFormatter string
	}

	// A SlogTextHandler is Logrus text formatter for log/slog.
	SlogTextHandler struct {
		opt    SlogTextOption
		writer io.Writer

		template string
		// Whether the logger's out is to a terminal.
		isTerminal bool
		// Color scheme to use.
		colorScheme *compiledColorScheme

		parent *SlogTextHandler
		prefix string
		group  string
		attrs  []slog.Attr
	}
)

// NewSlogTextHandler returns a new SlogTextHandler.
func NewSlogTextHandler(w io.Writer, o *SlogTextOption) *SlogTextHandler {
	if len(o.QuoteCharacter) == 0 {
		o.QuoteCharacter = "\""
	}
	if len(o.ValueFormatter) == 0 {
		o.ValueFormatter = "%v"
	}

	return &SlogTextHandler{
		writer:     w,
		opt:        *o,
		template:   " %s=" + o.ValueFormatter,
		isTerminal: checkIfTerminal(w),
	}
}

// Enabled reports whether the handler handles records at the given level.
func (h *SlogTextHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.opt.Level
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
func (h *SlogTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	nh := h.Clone()
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
func (h *SlogTextHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	nh := h.Clone()
	nh.group = name
	return nh
}

// Handle handles the Record.
func (h *SlogTextHandler) Handle(_ context.Context, record slog.Record) error {
	if h.prefix != "" {
		record.Message = fmt.Sprintf("%s %s", h.prefix, record.Message)
	}

	gprefix, keys, m := h.build()

	if record.NumAttrs() > 0 {
		record.Attrs(func(attr slog.Attr) bool {
			keys = append(keys, grouprecord(m, gprefix, attr)...)
			return true
		})
	}

	lastKeyIdx := len(keys) - 1

	if !h.opt.DisableSorting {
		sort.Strings(keys)
	}

	b := new(bytes.Buffer)
	prefixFieldClashes(m)

	timestampFormat := h.opt.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	if h.opt.ForceFormatting || h.isTerminal { // Is formatted
		isColored := (h.opt.ForceColors || h.isTerminal) && !h.opt.DisableColors
		var colorScheme *compiledColorScheme
		if isColored {
			colorScheme = h.colorScheme
			if colorScheme == nil {
				colorScheme = defaultCompiledColorScheme
			}
		} else {
			colorScheme = noColorsColorScheme
		}

		h.printColored(b, record, keys, m, timestampFormat, colorScheme)
	} else {
		h.appendKeyValue(b, "level", record.Level.String(), true)

		if !h.opt.DisableTimestamp {
			h.appendKeyValue(b, "time", record.Time.Format(timestampFormat), true)
		}

		if record.Message != "" {
			h.appendKeyValue(b, "msg", record.Message, lastKeyIdx >= 0)
		}

		for i, key := range keys {
			h.appendKeyValue(b, key, m[key], lastKeyIdx != i)
		}
	}

	b.WriteByte('\n')
	_, err := h.writer.Write(b.Bytes())
	return err
}

// Clone clones the entry, it creates a new instance and linking the parent to it.
func (h *SlogTextHandler) Clone() *SlogTextHandler {
	return &SlogTextHandler{
		parent:      h,
		template:    h.template,
		isTerminal:  h.isTerminal,
		colorScheme: h.colorScheme,
		prefix:      h.prefix,
		opt:         h.opt,
		writer:      h.writer,
	}
}

func (h *SlogTextHandler) build() (string, []string, map[string]any) {
	// Get all parents in a list.
	ilineage := make([]*SlogTextHandler, 0, 100)
	p := h
	for p != nil {
		ilineage = append(ilineage, p)
		p = p.parent
	}

	// Process attrs from parents to children.
	var gprefix string
	keys := make([]string, 0, 100)
	m := make(map[string]any)
	for i := len(ilineage) - 1; i >= 0; i-- {
		p = ilineage[i]
		if p.group != "" {
			gprefix += p.group + delimiter
		}

		for _, attr := range ilineage[i].attrs {
			k := gprefix + attr.Key
			if _, ok := m[k]; !ok {
				keys = append(keys, k)
			}

			m[k] = attr.Value.Any()
		}
	}

	return gprefix, keys, m
}

func (h *SlogTextHandler) appendKeyValue(b *bytes.Buffer, key string, value any, appendSpace bool) {
	b.WriteString(key)
	b.WriteByte('=')
	h.appendValue(b, value)

	if appendSpace {
		b.WriteByte(' ')
	}
}

func (h *SlogTextHandler) printColored(b *bytes.Buffer, record slog.Record, keys []string, m map[string]any, timestampFormat string, colorScheme *compiledColorScheme) {
	var levelColor func(string) string
	var levelText string
	switch record.Level {
	case slog.LevelInfo:
		levelColor = colorScheme.InfoLevelColor
	case slog.LevelWarn:
		levelColor = colorScheme.WarnLevelColor
	case slog.LevelError:
		levelColor = colorScheme.ErrorLevelColor
	default:
		levelColor = colorScheme.DebugLevelColor
	}

	levelText = "warn"
	if record.Level != slog.LevelWarn {
		levelText = record.Level.String()
	}

	if !h.opt.DisableUppercase {
		levelText = strings.ToUpper(levelText)
	}

	level := levelColor(fmt.Sprintf("%5s", levelText))

	message := record.Message
	if h.opt.PrefixRE != nil {
		if h.opt.PrefixRE.MatchString(message) {
			match := h.opt.PrefixRE.FindString(message)
			message = colorScheme.PrefixColor(message[:len(match)]) + message[len(match):]
		}
	}

	messageFormat := "%s"
	if h.opt.SpacePadding != 0 {
		messageFormat = fmt.Sprintf("%%-%ds", h.opt.SpacePadding)
	}

	if h.opt.DisableTimestamp {
		fmt.Fprintf(b, "%s "+messageFormat, level, message)
	} else {
		var timestamp string
		if !h.opt.FullTimestamp {
			timestamp = fmt.Sprintf("[%04d]", miniTS())
		} else {
			timestamp = fmt.Sprintf("[%s]", record.Time.Format(timestampFormat))
		}
		fmt.Fprintf(b, "%s %s "+messageFormat, colorScheme.TimestampColor(timestamp), level, message)
	}

	for _, k := range keys {
		fmt.Fprintf(b, h.template, levelColor(k), m[k])
	}
}

func (h *SlogTextHandler) appendValue(b *bytes.Buffer, value any) {
	switch value := value.(type) {
	case string:
		if !h.needsQuoting(value) {
			b.WriteString(value)
		} else {
			fmt.Fprintf(b, "%s%v%s", h.opt.QuoteCharacter, value, h.opt.QuoteCharacter)
		}
	case error:
		errmsg := value.Error()
		if !h.needsQuoting(errmsg) {
			b.WriteString(errmsg)
		} else {
			fmt.Fprintf(b, "%s%v%s", h.opt.QuoteCharacter, errmsg, h.opt.QuoteCharacter)
		}
	default:
		fmt.Fprint(b, value)
	}
}

func (h *SlogTextHandler) needsQuoting(text string) bool {
	if h.opt.QuoteEmptyFields && len(text) == 0 {
		return true
	}

	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.') {
			return true
		}
	}
	return false
}

func grouprecord(m map[string]any, p string, attr slog.Attr) []string {
	var keys []string

	if attr.Value.Kind() != slog.KindGroup {
		k := p + attr.Key
		if _, ok := m[k]; !ok {
			keys = append(keys, k)
		}

		m[k] = attr.Value.Any()
		return keys
	}

	p += attr.Key + delimiter
	for _, a := range attr.Value.Group() {
		if a.Value.Kind() == slog.KindGroup {
			keys = append(keys, grouprecord(m, p, a)...)
			continue
		}

		k := p + a.Key
		if _, ok := m[k]; !ok {
			keys = append(keys, k)
		}
		m[k] = a.Value.Any()
	}

	return keys
}
