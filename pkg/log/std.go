package log

import (
	"bytes"
	"io"
	"log"
	"strings"
	"sync"
	"text/template"
)

var _ Logger = (*stdLogger)(nil)

type stdLogger struct {
	log   *log.Logger
	level Level
	pool  *sync.Pool
	templ string
}

// NewStdLogger new a logger with writer.
func NewStdLogger(w io.Writer, level Level) Logger {
	l := &stdLogger{
		log:   log.New(w, "", 0),
		level: level,
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		templ: "{{.Timestamp}} {{.Level}} {{.Caller}} {{.Msg}}",
	}

	return With(l, "timestamp", DefaultTimestamp, "caller", DefaultCaller)
}

func NewStdLoggerWithFormat(w io.Writer, level Level, fmt string) Logger {
	l := &stdLogger{
		log:   log.New(w, "", 0),
		level: level,
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		templ: fmt,
	}

	return With(l, "timestamp", DefaultTimestamp, "caller", DefaultCaller)
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// Log print the kv pairs log.
func (l *stdLogger) Log(level Level, keyvals ...interface{}) error {
	if level < l.level {
		return nil
	}
	if len(keyvals) == 0 {
		return nil
	}
	if (len(keyvals) & 1) == 1 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
	}
	buf := l.pool.Get().(*bytes.Buffer)

	values := make(map[string]any)
	values["Level"] = level.String()

	for i := 0; i < len(keyvals); i += 2 {
		if key, ok := keyvals[i].(string); ok {
			values[capitalize(key)] = keyvals[i+1]
		}
	}
	tmpl, err := template.New("log").Parse(l.templ)
	if err != nil {
		return err
	}
	err = tmpl.Execute(buf, values)
	if err != nil {
		return err
	}

	_ = l.log.Output(depth, buf.String()) //nolint:gomnd
	buf.Reset()
	l.pool.Put(buf)
	return nil
}

func (l *stdLogger) Close() error {
	return nil
}
