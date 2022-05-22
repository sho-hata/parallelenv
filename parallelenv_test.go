package parallelenv

import (
	"testing"

	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, Analyzer, "a")
}

func Test_target_isValid(t *testing.T) {
	t.Parallel()
	isValid, isNotValid := true, false
	tests := map[string]struct {
		tr   target
		want bool
	}{
		"valid value":        {tSetEnv, isValid},
		"unknown value":      {unknown, isNotValid},
		"out of range value": {-1, isNotValid},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := tt.tr.isValid(); got != tt.want {
				t.Errorf("target.isValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sToTarget(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		s    string
		want target
	}{
		"success": {"t.Setenv", tSetEnv},
		"failed":  {"t.Short", unknown},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := sToTarget(tt.s); got != tt.want {
				t.Errorf("sToTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}
