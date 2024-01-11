package logger

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

// A LogrusTextFormatter is pretty printer for Logrus.
// It supports github.com/mdouchement/logger prefix.
// Borrowed from https://github.com/x-cray/logrus-prefixed-formatter.
type LogrusTextFormatter struct {
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

	// Color scheme to use.
	colorScheme *compiledColorScheme

	// Whether the logger's out is to a terminal.
	isTerminal bool

	// Index of the printed entry (useful to sort lines in Graylog)
	index uint32

	template string

	sync.Once
}

// SetColorScheme setup the color scheme.
func (f *LogrusTextFormatter) SetColorScheme(colorScheme *ColorScheme) {
	f.colorScheme = compileColorScheme(colorScheme)
}

// Format implements logrus.Formatter.
func (f *LogrusTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if prefix, ok := entry.Data[KeyPrefix]; ok {
		entry.Message = fmt.Sprintf("%s %s", prefix, entry.Message)
		delete(entry.Data, KeyPrefix)
	}

	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	lastKeyIdx := len(keys) - 1

	if !f.DisableSorting {
		sort.Strings(keys)
	}

	b := entry.Buffer
	if b != nil {
		b = new(bytes.Buffer)
	}

	prefixFieldClashes(entry.Data)

	f.Do(func() { f.init(entry) })

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	if f.ForceFormatting || f.isTerminal { // Is formatted
		isColored := (f.ForceColors || f.isTerminal) && !f.DisableColors
		var colorScheme *compiledColorScheme
		if isColored {
			colorScheme = f.colorScheme
			if colorScheme == nil {
				colorScheme = defaultCompiledColorScheme
			}
		} else {
			colorScheme = noColorsColorScheme
		}
		f.printColored(b, entry, keys, timestampFormat, colorScheme)
	} else {
		index := atomic.AddUint32(&f.index, 1)
		f.appendKeyValue(b, "index", index, true)
		f.appendKeyValue(b, "level", entry.Level.String(), true)

		if !f.DisableTimestamp {
			f.appendKeyValue(b, "time", entry.Time.Format(timestampFormat), true)
		}

		if entry.Message != "" {
			f.appendKeyValue(b, "msg", entry.Message, lastKeyIdx >= 0)
		}

		for i, key := range keys {
			f.appendKeyValue(b, key, entry.Data[key], lastKeyIdx != i)
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *LogrusTextFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry, keys []string, timestampFormat string, colorScheme *compiledColorScheme) {
	var levelColor func(string) string
	var levelText string
	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = colorScheme.InfoLevelColor
	case logrus.WarnLevel:
		levelColor = colorScheme.WarnLevelColor
	case logrus.ErrorLevel:
		levelColor = colorScheme.ErrorLevelColor
	case logrus.FatalLevel:
		levelColor = colorScheme.FatalLevelColor
	case logrus.PanicLevel:
		levelColor = colorScheme.PanicLevelColor
	default:
		levelColor = colorScheme.DebugLevelColor
	}

	levelText = "warn"
	if entry.Level != logrus.WarnLevel {
		levelText = entry.Level.String()
	}

	if !f.DisableUppercase {
		levelText = strings.ToUpper(levelText)
	}

	level := levelColor(fmt.Sprintf("%5s", levelText))

	message := entry.Message
	if f.PrefixRE != nil {
		if f.PrefixRE.MatchString(message) {
			match := f.PrefixRE.FindString(message)
			message = colorScheme.PrefixColor(message[:len(match)]) + message[len(match):]
		}
	}

	messageFormat := "%s"
	if f.SpacePadding != 0 {
		messageFormat = fmt.Sprintf("%%-%ds", f.SpacePadding)
	}

	if f.DisableTimestamp {
		fmt.Fprintf(b, "%s "+messageFormat, level, message)
	} else {
		var timestamp string
		if !f.FullTimestamp {
			timestamp = fmt.Sprintf("[%04d]", miniTS())
		} else {
			timestamp = fmt.Sprintf("[%s]", entry.Time.Format(timestampFormat))
		}
		fmt.Fprintf(b, "%s %s "+messageFormat, colorScheme.TimestampColor(timestamp), level, message)
	}

	for _, k := range keys {
		v := entry.Data[k]
		fmt.Fprintf(b, f.template, levelColor(k), v)
	}
}

func (f *LogrusTextFormatter) needsQuoting(text string) bool {
	if f.QuoteEmptyFields && len(text) == 0 {
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

func (f *LogrusTextFormatter) appendKeyValue(b *bytes.Buffer, key string, value any, appendSpace bool) {
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)

	if appendSpace {
		b.WriteByte(' ')
	}
}

func (f *LogrusTextFormatter) appendValue(b *bytes.Buffer, value any) {
	switch value := value.(type) {
	case string:
		if !f.needsQuoting(value) {
			b.WriteString(value)
		} else {
			fmt.Fprintf(b, "%s%v%s", f.QuoteCharacter, value, f.QuoteCharacter)
		}
	case error:
		errmsg := value.Error()
		if !f.needsQuoting(errmsg) {
			b.WriteString(errmsg)
		} else {
			fmt.Fprintf(b, "%s%v%s", f.QuoteCharacter, errmsg, f.QuoteCharacter)
		}
	default:
		fmt.Fprint(b, value)
	}
}

func (f *LogrusTextFormatter) init(entry *logrus.Entry) {
	if len(f.QuoteCharacter) == 0 {
		f.QuoteCharacter = "\""
	}

	if len(f.ValueFormatter) == 0 {
		f.ValueFormatter = "%v"
	}
	f.template = " %s=" + f.ValueFormatter

	if entry.Logger != nil {
		f.isTerminal = checkIfTerminal(entry.Logger.Out)
	}
}
