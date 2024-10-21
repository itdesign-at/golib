package xlog

import (
	"log/slog"
	"testing"
)

func TestParseLevel(t *testing.T) {
	t.Parallel()

	tests := map[string]slog.Level{
		"debug":   slog.LevelDebug,
		"info":    slog.LevelInfo,
		"audit":   LevelAudit,
		"warn":    slog.LevelWarn,
		"error":   slog.LevelError,
		"invalid": LevelAudit,
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			t.Parallel()

			o := ClientOptions{
				Level: input,
			}

			l := o.parseLevel()
			if l != expected {
				t.Errorf("expected %v, got %v", expected, l)
			}
		})
	}

}
