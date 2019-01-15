package config

import (
	"testing"
)

func TestEndpointToKey(t *testing.T) {
	type fields struct {
		Endpoint string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"normal", fields{"https://api.altizure.com/graphql"}, "https://api*altizure*com"},
		{"empty", fields{""}, "https://api*altizure*com"},
		{"invalid", fields{"nat"}, "https://api*altizure*com"},
		{"no scheme", fields{"nat.com"}, "https://api*altizure*com"},
		{"with port", fields{"http://nat.com:12345/abcde"}, "http://nat*com:12345"},
		{"localhost", fields{"http://127.0.0.1:8082"}, "http://127*0*0*1:8082"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := endpointToKey(tt.fields.Endpoint); got != tt.want {
				t.Errorf("Scope.Key() = %v, want %v", got, tt.want)
			}
		})
	}
}
