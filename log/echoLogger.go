package log

import (
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"time"
	"bytes"
	"runtime"
	"strconv"
	"encoding/json"

	"github.com/mattn/go-isatty"
	"github.com/mattn/go-colorable"
	"github.com/labstack/gommon/log"
	"github.com/valyala/fasttemplate"

	clr "github.com/labstack/gommon/color"
)

type (
	EchoLogger struct {
		prefix     string
		level      log.Lvl
		output     io.Writer
		template   *fasttemplate.Template
		levels     []string
		color      *clr.Color
		bufferPool sync.Pool
		mutex      sync.Mutex
	}

)

const (
	Gray = uint8(iota + 90)
	Red
	Green
	Yellow
	Blue
	Magenta
	EndColor = "\033[0m"

	INFO = "INFO"
	TRAC = "TRAC"
	ERRO = "ERRO"
	WARN = "WARN"
	SUCC = "SUCC"
)

var (
	global = New("-")
	defaultHeader = "${time_rfc3339} [${level}]"
	Color = clr.New()
)

func New(prefix string) (l *EchoLogger) {
	l = &EchoLogger{
		level:    log.INFO,
		prefix:   prefix,
		template: l.newTemplate(defaultHeader),
		color:    Color,
		bufferPool: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 256))
			},
		},
	}

	l.initLevels()
	l.SetOutput(colorable.NewColorableStdout())
	return
}

func (l *EchoLogger) initLevels() {
	l.levels = []string{
		l.color.Blue("DEBUG"),
		l.color.Green("INFO"),
		l.color.Yellow("WARN"),
		l.color.Red("ERROR"),
		l.color.RedBg("FATAL"),
	}
}

func (l *EchoLogger) newTemplate(format string) *fasttemplate.Template {
	return fasttemplate.New(format, "${", "}")
}

func (l *EchoLogger) DisableColor() {
	l.color.Disable()
	l.initLevels()
}

func (l *EchoLogger) EnableColor() {
	l.color.Enable()
	l.initLevels()
}

func (l *EchoLogger) Prefix() string {
	return l.prefix
}

func (l *EchoLogger) SetPrefix(p string) {
	l.prefix = p
}

func (l *EchoLogger) Level() log.Lvl {
	return l.level
}

func (l *EchoLogger) SetLevel(v log.Lvl) {
	l.level = v
}

func (l *EchoLogger) Output() io.Writer {
	return l.output
}

func (l *EchoLogger) SetHeader(h string) {
	l.template = l.newTemplate(h)
}

func (l *EchoLogger) SetOutput(w io.Writer) {
	l.output = w
	if w, ok := w.(*os.File); !ok || !isatty.IsTerminal(w.Fd()) {
		l.DisableColor()
	}
}

func (l *EchoLogger) Print(i ...interface{}) {
	fmt.Fprintln(l.output, i...)
}

func (l *EchoLogger) Printf(format string, args ...interface{}) {
	f := fmt.Sprintf("%s\n", format)
	fmt.Fprintf(l.output, f, args...)
}

func (l *EchoLogger) Printj(j log.JSON) {
	json.NewEncoder(l.output).Encode(j)
}

func (l *EchoLogger) Debug(i ...interface{}) {
	l.log(log.DEBUG, "", i...)
}

func (l *EchoLogger) Debugf(format string, args ...interface{}) {
	l.log(log.DEBUG, format, args...)
}

func (l *EchoLogger) Debugj(j log.JSON) {
	l.log(log.DEBUG, "json", j)
}

func (l *EchoLogger) Info(i ...interface{}) {
	l.log(log.INFO, "", i...)
}

func (l *EchoLogger) Infof(format string, args ...interface{}) {
	l.log(log.INFO, format, args...)
}

func (l *EchoLogger) Infoj(j log.JSON) {
	l.log(log.INFO, "json", j)
}

func (l *EchoLogger) Warn(i ...interface{}) {
	l.log(log.WARN, "", i...)
}

func (l *EchoLogger) Warnf(format string, args ...interface{}) {
	l.log(log.WARN, format, args...)
}

func (l *EchoLogger) Warnj(j log.JSON) {
	l.log(log.WARN, "json", j)
}

func (l *EchoLogger) Error(i ...interface{}) {
	l.log(log.ERROR, "", i...)
}

func (l *EchoLogger) Errorf(format string, args ...interface{}) {
	l.log(log.ERROR, format, args...)
}

func (l *EchoLogger) Errorj(j log.JSON) {
	l.log(log.ERROR, "json", j)
}

func (l *EchoLogger) Fatal(i ...interface{}) {
	l.log(log.FATAL, "", i...)
	os.Exit(1)
}

func (l *EchoLogger) Fatalf(format string, args ...interface{}) {
	l.log(log.FATAL, format, args...)
	os.Exit(1)
}

func (l *EchoLogger) Fatalj(j log.JSON) {
	l.log(log.FATAL, "json", j)
}

func DisableColor() {
	global.DisableColor()
}

func EnableColor() {
	global.EnableColor()
}

func Prefix() string {
	return global.Prefix()
}

func SetPrefix(p string) {
	global.SetPrefix(p)
}

func Level() log.Lvl {
	return global.Level()
}

func SetLevel(v log.Lvl) {
	global.SetLevel(v)
}

func Output() io.Writer {
	return global.Output()
}

func SetOutput(w io.Writer) {
	global.SetOutput(w)
}

func SetHeader(h string) {
	global.SetHeader(h)
}

func Print(i ...interface{}) {
	global.Print(i...)
}

func Printf(format string, args ...interface{}) {
	global.Printf(format, args...)
}

func Printj(j log.JSON) {
	global.Printj(j)
}

func Debug(i ...interface{}) {
	global.Debug(i...)
}

func Debugf(format string, args ...interface{}) {
	global.Debugf(format, args...)
}

func Debugj(j log.JSON) {
	global.Debugj(j)
}

func Info(i ...interface{}) {
	global.Info(i...)
}

func Infof(format string, args ...interface{}) {
	global.Infof(format, args...)
}

func Infoj(j log.JSON) {
	global.Infoj(j)
}

func Warn(i ...interface{}) {
	global.Warn(i...)
}

func Warnf(format string, args ...interface{}) {
	global.Warnf(format, args...)
}

func Warnj(j log.JSON) {
	global.Warnj(j)
}

func Error(i ...interface{}) {
	global.Error(i...)
}

func Errorf(format string, args ...interface{}) {
	global.Errorf(format, args...)
}

func Errorj(j log.JSON) {
	global.Errorj(j)
}

func Fatal(i ...interface{}) {
	global.Fatal(i...)
}

func Fatalf(format string, args ...interface{}) {
	global.Fatalf(format, args...)
}

func Fatalj(j log.JSON) {
	global.Fatalj(j)
}

func (l *EchoLogger) log(v log.Lvl, format string, args ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	buf := l.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer l.bufferPool.Put(buf)
	_, file, line, _ := runtime.Caller(3)

	if v >= l.level {
		message := ""
		if format == "" {
			message = fmt.Sprint(args...)
		} else if format == "json" {
			b, err := json.Marshal(args[0])
			if err != nil {
				panic(err)
			}
			message = string(b)
		} else {
			message = fmt.Sprintf(format, args...)
		}

		if v >= log.ERROR {
			// panic(message)
		}

		_, err := l.template.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
			switch tag {
			case "time_rfc3339":
				return w.Write([]byte(time.Now().Format("2006/01/02 15:04:05")))
			case "level":
				return w.Write([]byte(l.levels[v]))
			case "prefix":
				return w.Write([]byte(l.prefix))
			case "long_file":
				return w.Write([]byte(file))
			case "short_file":
				return w.Write([]byte(path.Base(file)))
			case "line":
				return w.Write([]byte(strconv.Itoa(line)))
			}
			return 0, nil
		})

		if err == nil {
			s := buf.String()
			i := buf.Len() - 1
			if s[i] == '}' {
				// log.JSON header
				buf.Truncate(i)
				buf.WriteByte(',')
				if format == "json" {
					buf.WriteString(message[1:])
				} else {
					buf.WriteString(`"message":"`)
					buf.WriteString(message)
					buf.WriteString(`"}`)
				}
			} else {
				// Text header
				buf.WriteByte(' ')
				buf.WriteString(message)
			}
			buf.WriteByte('\n')
			l.output.Write(buf.Bytes())
		}
	}
}
