package clock

import (
	"reflect"
	"testing"
	"time"
)

func TestBeginningOfDay(t *testing.T) {
	t.Parallel()

	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			args: args{t: time.Date(2020, 1, 1, 23, 59, 59, 0, time.UTC)},
			want: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			args: args{t: time.Date(2020, 2, 28, 0, 1, 0, 0, time.UTC)},
			want: time.Date(2020, 2, 28, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := BeginningOfDay(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BeginningOfDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndOfDay(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			args: args{t: time.Date(2020, 1, 1, 23, 59, 58, 0, time.UTC)},
			want: time.Date(2020, 1, 1, 23, 59, 59, 999999999, time.UTC),
		},
		{
			args: args{t: time.Date(2020, 2, 28, 0, 1, 0, 0, time.UTC)},
			want: time.Date(2020, 2, 28, 23, 59, 59, 999999999, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndOfDay(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EndOfDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetTimeZone(t *testing.T) {
	type args struct {
		timeZone string
	}
	tests := []struct {
		name    string
		args    args
		want string
		wantErr bool
	}{
		{
			name: "invalid timezone",
			args: args{
				timeZone: "Invalid",
			},
			want: "Local",
			wantErr: true,
		},
		{
			name: "valid timezone",
			args: args{
				timeZone: "Asia/Tokyo",
			},
			want: "Asia/Tokyo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetTimeZone(tt.args.timeZone); (err != nil) != tt.wantErr {
				t.Errorf("SetTimeZone() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := time.Local.String(); got != tt.want {
				t.Errorf("SetTimeZone() got = %v, want = %v", got, tt.want)
			}
		})
	}
}
