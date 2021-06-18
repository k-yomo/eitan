package appenv

import "testing"

func TestEnv_IsDeployed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		e    Env
		want bool
	}{
		{
			e:    Test,
			want: false,
		},
		{
			e:    Local,
			want: false,
		},
		{
			e:    Dev,
			want: true,
		},
		{
			e:    Prod,
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.e.IsDeployed(); got != tt.want {
				t.Errorf("IsDeployed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    Env
		want bool
	}{
		{
			e:    Env("invalid"),
			want: false,
		},
		{
			e:    Prod,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
