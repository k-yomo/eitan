package randnum

import (
	"regexp"
	"testing"
)

func TestRandNumString(t *testing.T) {
	t.Parallel()

	type args struct {
		length int
	}
	tests := []struct {
		name      string
		args      args
		wantRegex *regexp.Regexp
	}{
		{
			args:      args{3},
			wantRegex: regexp.MustCompile("[0-9]{3}"),
		},
		{
			args:      args{10},
			wantRegex: regexp.MustCompile("[0-9]{10}"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := RandNumString(tt.args.length); !tt.wantRegex.Match([]byte(got)) {
				t.Errorf("RandNumString() = %v, wantRegex %v", got, tt.wantRegex.String())
			}
		})
	}
}
