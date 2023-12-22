package logger_test

import (
	"bytes"
	"log/slog"
	"regexp"
	"testing"
	"time"

	"github.com/mdouchement/logger"
)

func TestSlog(t *testing.T) {
	w := new(bytes.Buffer)
	l := slog.New(logger.NewSlogGELFHandler(w, &logger.SlogGELFOption{Level: slog.LevelInfo, Hostname: "hostname-42"}))

	l = l.With(logger.KeyPrefix, "[uuid]")
	l.With(logger.KeyPrefix, "[void]") // must not appear in message
	l = l.With("root", "root")

	l = l.WithGroup("g1")
	l = l.With("r1", "1")

	l = l.WithGroup("g2").WithGroup("g3")
	l = l.With("r", "x")
	l = l.With("l", "y")

	l.Debug("debug")
	l.Info("info")
	l.Warn("warn")
	l.Error("error")

	// INFO
	expected := regexp.MustCompile(`\{"_g1\.g2\.g3\.l":"y","_g1\.g2\.g3\.r":"x","_g1\.r1":"1","_level_name":"INFO","_root":"root","host":"hostname-42","level":6,"short_message":"\[uuid\] info","timestamp":\d+.\d+,"version":"1\.1"\}`)
	line, err := w.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if !expected.MatchString(line) {
		t.Errorf("got: %s", line)
	}

	// WARN
	expected = regexp.MustCompile(`\{"_g1\.g2\.g3\.l":"y","_g1\.g2\.g3\.r":"x","_g1\.r1":"1","_level_name":"WARN","_root":"root","host":"hostname-42","level":4,"short_message":"\[uuid\] warn","timestamp":\d+.\d+,"version":"1\.1"\}`)
	line, err = w.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if !expected.MatchString(line) {
		t.Errorf("got: %s", line)
	}

	// ERROR
	expected = regexp.MustCompile(`\{"_g1\.g2\.g3\.l":"y","_g1\.g2\.g3\.r":"x","_g1\.r1":"1","_level_name":"ERROR","_root":"root","host":"hostname-42","level":3,"short_message":"\[uuid\] error","timestamp":\d+.\d+,"version":"1\.1"\}`)
	line, err = w.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if !expected.MatchString(line) {
		t.Errorf("got: %s", line)
	}
}

func TestSlogRecord(t *testing.T) {
	w := new(bytes.Buffer)
	l := slog.New(logger.NewSlogGELFHandler(w, &logger.SlogGELFOption{Level: slog.LevelDebug, Hostname: "hostname-42"}))

	l = l.With(logger.KeyPrefix, "[uuid]")
	l = l.WithGroup("g1")
	l = l.With("r1", "1")

	l.Debug("debug", slog.Int("rec", 42), slog.Group("grec", slog.Bool("b", true), slog.Duration("t", time.Second)))

	expected := regexp.MustCompile(`\{"_g1\.grec\.b":"true","_g1\.grec\.t":"1s","_g1\.r1":"1","_g1\.rec":42,"_level_name":"DEBUG","host":"hostname-42","level":7,"short_message":"\[uuid\] debug","timestamp":\d+.\d+,"version":"1\.1"\}`)
	line, err := w.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if !expected.MatchString(line) {
		t.Errorf("got: %s", line)
	}
}
