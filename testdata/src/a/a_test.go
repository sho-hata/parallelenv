package a

import "testing"

func Test_f(t *testing.T) {
	t.Skip()
	t.Parallel()                 // want "cannot set environment variables in parallel tests"
	t.Setenv("SAMPLE", "sample") // want "cannot set environment variables in parallel tests"

	tests := []struct {
		name    string
		s       string
		want    string
		wantErr bool
	}{}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := f(tt.s)
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

func Test_f_1(t *testing.T) {
	t.Skip()
	t.Parallel() // want "cannot set environment variables in parallel tests"
	flag := true
	if flag {
		t.Setenv("SAMPLE", "sample") // want "cannot set environment variables in parallel tests"
	} else {
		t.Setenv("SAMPLE", "sample") // want "cannot set environment variables in parallel tests"
	}

	tests := []struct {
		name    string
		s       string
		want    string
		wantErr bool
	}{}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := f(tt.s)
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

func Test_f_2(t *testing.T) {
	t.Skip()

	tests := []struct {
		name    string
		s       string
		want    string
		wantErr bool
	}{}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()                 // want "cannot set environment variables in parallel tests"
			t.Setenv("SAMPLE", "sample") // want "cannot set environment variables in parallel tests"
			got, err := f(tt.s)
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

func Test_f_3(t *testing.T) {
	t.Skip()
	t.Parallel()

	tests := []struct {
		name    string
		s       string
		want    string
		wantErr bool
	}{}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SAMPLE", "sample")
			got, err := f(tt.s)
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

func Test_f_4(t *testing.T) {
	t.Skip()
	t.Parallel()                 // want "cannot set environment variables in parallel tests"
	t.Setenv("SAMPLE", "sample") // want "cannot set environment variables in parallel tests"

	tests := []struct {
		name    string
		s       string
		want    string
		wantErr bool
	}{}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := f(tt.s)
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
