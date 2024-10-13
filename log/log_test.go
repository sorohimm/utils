package log

import (
	"testing"

	"go.uber.org/zap"
)

func TestNewZap(t *testing.T) {
	testCases := []struct {
		name     string
		lvl      string
		encType  string
		wantErr  bool
		wantLvl  zap.AtomicLevel
		wantEnc  string
		wantCall bool
	}{
		{
			name:    "valid level and invalid encoding",
			lvl:     "debug",
			encType: "invalid",
			wantErr: true,
		},
		{
			name:    "invalid level and valid encoding",
			lvl:     "invalid",
			encType: "json",
			wantErr: true,
		},
		{
			name:    "empty level and valid encoding",
			lvl:     "",
			encType: "json",
			wantLvl: zap.NewAtomicLevelAt(zap.InfoLevel),
			wantEnc: "json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger, err := NewZap(tc.lvl, tc.encType)
			if (err != nil) != tc.wantErr {
				t.Errorf("NewZap() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if err != nil {
				return
			}

			if got := logger.Core().Enabled(zap.DebugLevel); got != tc.wantCall {
				t.Errorf("NewZap() enabled = %v, want %v", got, tc.wantCall)
			}
		})
	}
}
