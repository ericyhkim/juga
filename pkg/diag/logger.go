package diag

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

var (
	styleError = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444"))
	styleWarn  = lipgloss.NewStyle().Foreground(lipgloss.Color("#EAB308"))
	styleDebug = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
)

type StdLogger struct{}

func NewStdLogger() *StdLogger {
	return &StdLogger{}
}

func (l *StdLogger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fmt.Fprintln(os.Stderr, styleError.Render(msg))
}

func (l *StdLogger) Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fmt.Fprintln(os.Stderr, styleWarn.Render(msg))
}

func (l *StdLogger) Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fmt.Fprintln(os.Stderr, styleDebug.Render(msg))
}

type NopLogger struct{}

func NewNopLogger() *NopLogger {
	return &NopLogger{}
}

func (l *NopLogger) Error(format string, v ...interface{}) {}
func (l *NopLogger) Warn(format string, v ...interface{})  {}
func (l *NopLogger) Debug(format string, v ...interface{}) {}
