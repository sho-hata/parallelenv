package a

import "testing"

func Test_f(t *testing.T) {
	t.Skip()
	t.Parallel()                 // want "cannot set environment variables in parallel tests"
	t.Setenv("SAMPLE", "sample") // want "cannot set environment variables in parallel tests"

	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := f(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("f() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("f() = %v, want %v", got, tt.want)
			}
		})
	}
}
