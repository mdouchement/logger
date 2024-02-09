package logger_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/mdouchement/logger"
)

func TestBufferGELF(t *testing.T) {
	b := logger.NewBufferGELF()
	if string(b.Complete(false)) != `{"version":"1.1"}` {
		t.Errorf("got: %s", b.Bytes())
	}

	//

	b = logger.NewBufferGELF()
	if string(b.Complete(true)) != "{\"version\":\"1.1\"}\n" {
		t.Errorf("got: %s", b.Bytes())
	}

	//

	now := time.Now().Truncate(time.Nanosecond)
	expected := fmt.Sprintf(
		`{"version":"1.1","timestamp":%s,"host":"host","level":7,"short_message":"message","_k1":42,"_k2":42.42,"_k3":42,"_k4":"true","_k5":%s}`,
		strconv.FormatFloat(float64(now.UnixNano())/1e9, 'f', -1, 64),
		strconv.QuoteToGraphic(now.Format(time.RFC3339)),
	)
	b = logger.NewBufferGELF()
	b.Timestamp(now)
	b.Host("host")
	b.Level(7)
	b.Message("message")
	b.Add("k1", 42)
	b.Add("k2", 42.42)
	b.Add("k3", uint64(42))
	b.Add("k4", true)
	b.Add("k5", now)
	if string(b.Complete(false)) != expected {
		t.Errorf("\n   got: %s\nexpect: %s", b.Bytes(), expected)
	}

	//

	b = logger.NewBufferGELF()
	b.Message("first\ntwo")
	if string(b.Complete(false)) != `{"version":"1.1","short_message":"first","full_message":"first\ntwo"}` {
		t.Errorf("got: %s", b.Bytes())
	}

	//

	b = logger.NewBufferGELF()
	b.Add("k", 42)
	b.Add("k", 43)
	if string(b.Complete(false)) != `{"version":"1.1","_k":42,"_k":43}` {
		t.Errorf("got: %s", b.Bytes())
	}

	//

	b = logger.NewBufferGELF()
	b.Message(string([]byte{162, 98, 107, 50, 251, 64, 64, 184, 81, 235, 133, 30, 184, 99, 107, 101, 121, 24, 42})) // Some CBOR payload
	if string(b.Complete(false)) != `{"version":"1.1","short_message":"\ufffdbk2\ufffd@@\ufffdQ\ufffd\ufffd\u001e\ufffdckey\u0018*"}` {
		t.Errorf("got: %s", b.Bytes())
	}
}
