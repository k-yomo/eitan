package csrf

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewCSRFValidationMiddleware(t *testing.T) {
	t.Parallel()
	type args struct {
		next       http.Handler
		csrfHeader string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Continues to handle request when X-Requested-By is set",
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
				csrfHeader: "some value",
			},
			want: 200,
		},
		{
			name: "Returns 422 when X-Requested-By is empty",
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
				csrfHeader: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "http://testing", nil)
			req.Header.Add("X-Requested-By", tt.args.csrfHeader)
			handlerToTest := NewCSRFValidationMiddleware(true)(tt.args.next)
			recorder := httptest.NewRecorder()
			handlerToTest.ServeHTTP(recorder, req)
		})
	}
}
