package logger

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Logger はログ出力を管理する構造体
type Logger struct {
	debug bool
}

// New は新しいロガーを生成します
func New() *Logger {
	return &Logger{
		debug: os.Getenv("DEBUG") == "1",
	}
}

// Info は情報メッセージを出力します
func (l *Logger) Info(msg string) {
	if _, err := fmt.Fprintln(os.Stdout, msg); err != nil {
		color.Red("log output error: %v", err)
	}
}

// Infof はフォーマット付きの情報メッセージを出力します
func (l *Logger) Infof(format string, args ...interface{}) {
	if _, err := fmt.Fprintf(os.Stdout, format+"\n", args...); err != nil {
		color.Red("log output error: %v", err)
	}
}

// Error はエラーメッセージを出力します
func (l *Logger) Error(msg string) {
	color.Red(msg)
}

// Errorf はフォーマット付きのエラーメッセージを出力します
func (l *Logger) Errorf(format string, args ...interface{}) {
	color.Red(format, args...)
}

// Success は成功メッセージを出力します
func (l *Logger) Success(msg string) {
	color.Green(msg)
}

// Successf はフォーマット付きの成功メッセージを出力します
func (l *Logger) Successf(format string, args ...interface{}) {
	color.Green(format, args...)
}

// Debug はデバッグメッセージを出力します（DEBUGフラグが有効な場合のみ）
func (l *Logger) Debug(msg string) {
	if l.debug {
		color.Yellow("[DEBUG] %s", msg)
	}
}

// Debugf はフォーマット付きのデバッグメッセージを出力します（DEBUGフラグが有効な場合のみ）
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.debug {
		color.Yellow("[DEBUG] "+format, args...)
	}
}
