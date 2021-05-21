package customerror

import (
	"github.com/pkg/errors"
	"testing"
)

func Test_customError_Error(t *testing.T) {
	t.Parallel()
	originalErr := errors.New("some error")
	type fields struct {
		errType ErrType
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "return original error's Error",
			fields: fields{err: originalErr},
			want:   originalErr.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ce := &customError{
				ErrType: tt.fields.errType,
				err:     tt.fields.err,
			}
			if got := ce.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
