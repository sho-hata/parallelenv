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
		"valid value":        {setEnv, isValid},
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
		"success": {"t.Setenv", setEnv},
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

func Test_checkState_shouldReport(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		s    checkState
		want bool
	}{
		"should report":                {checkState{parallel: {11111}, setEnv: {111111}}, true},
		"only the Setenv was called":   {checkState{setEnv: {11111}}, false},
		"only the parallel was called": {checkState{parallel: {11111}}, false},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := tt.s.shouldReport(); got != tt.want {
				t.Errorf("checkState.shouldReport() = %v, want %v", got, tt.want)
			}
		})
	}
}
