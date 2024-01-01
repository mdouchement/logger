package logger_test

import (
	"bytes"
	"log/slog"
	"regexp"
	"testing"

	"github.com/mdouchement/logger"
	"github.com/sirupsen/logrus"
)

// go test -bench . -run NONE -benchtime 5s

var result []byte

func BenchmarkLogrusGELF(b *testing.B) {
	w := new(bytes.Buffer)
	ll := logrus.New()
	ll.SetOutput(w)
	ll.SetFormatter(&logger.LogrusGELFFormatter{Hostname: "hostname"})
	base := logger.WrapLogrus(ll)

	var l logger.Logger
	for i := 0; i < b.N; i++ {
		l = base.WithPrefix("[prefix]")
		l = l.WithField("f1", 42).WithField("f2", "42").WithField("f3", 42).WithField("f4", "42")
		l = l.WithPrefixf("[%s]", 4242)
		l.Info("message")
	}

	b.ReportAllocs()
	result = w.Bytes()
}

func BenchmarkSlogGELF(b *testing.B) {
	w := new(bytes.Buffer)
	base := logger.NewPrefixLogger(&logger.PrefixGELF{Writer: w, Level: slog.LevelInfo, Hostname: "hostname"})

	var l logger.Logger
	for i := 0; i < b.N; i++ {
		l = base.WithPrefix("[prefix]")
		l = l.WithField("f1", 42).WithField("f2", "42").WithField("f3", 42).WithField("f4", "42")
		l = l.WithPrefixf("[%s]", 4242)
		l.Info("message")
	}

	b.ReportAllocs()
	result = w.Bytes()
}

func BenchmarkPrefixLoggerGELF(b *testing.B) {
	w := new(bytes.Buffer)
	ll := slog.New(logger.NewSlogGELFHandler(w, &logger.SlogGELFOption{Hostname: "hostname"}))
	base := logger.WrapSlog(ll)

	var l logger.Logger
	for i := 0; i < b.N; i++ {
		l = base.WithPrefix("[prefix]")
		l = l.WithField("f1", 42).WithField("f2", "42").WithField("f3", 42).WithField("f4", "42")
		l = l.WithPrefixf("[%s]", 4242)
		l.Info("message")
	}

	b.ReportAllocs()
	result = w.Bytes()
}

func BenchmarkLogrusText(b *testing.B) {
	w := new(bytes.Buffer)
	ll := logrus.New()
	ll.SetOutput(w)
	ll.SetFormatter(&logger.LogrusTextFormatter{
		DisableColors:   false,
		ForceColors:     true,
		ForceFormatting: true,
		PrefixRE:        regexp.MustCompile(`^(\[.*?\])\s`),
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	base := logger.WrapLogrus(ll)

	var l logger.Logger
	for i := 0; i < b.N; i++ {
		l = base.WithPrefix("[prefix]")
		l = l.WithField("f1", 42).WithField("f2", "42").WithField("f3", 42).WithField("f4", "42")
		l = l.WithPrefixf("[%s]", 4242)
		l.Info("message")
	}

	b.ReportAllocs()
	result = w.Bytes()
}

func BenchmarkSlogText(b *testing.B) {
	w := new(bytes.Buffer)
	ll := slog.New(logger.NewSlogTextHandler(w, &logger.SlogTextOption{
		DisableColors:   false,
		ForceColors:     true,
		ForceFormatting: true,
		PrefixRE:        regexp.MustCompile(`^(\[.*?\])\s`),
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}))
	base := logger.WrapSlog(ll)

	var l logger.Logger
	for i := 0; i < b.N; i++ {
		l = base.WithPrefix("[prefix]")
		l = l.WithField("f1", 42).WithField("f2", "42").WithField("f3", 42).WithField("f4", "42")
		l = l.WithPrefixf("[%s]", 4242)
		l.Info("message")
	}

	b.ReportAllocs()
	result = w.Bytes()
}
