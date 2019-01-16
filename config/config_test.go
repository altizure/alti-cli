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

func TestUniqueProfile(t *testing.T) {
	ps := []Profile{
		{ID: "id1", Key: "k1", Token: "t1"},
		{ID: "id2", Key: "k2", Token: "t2"},
		{ID: "id3", Key: "k1", Token: "t1"},
	}
	ups := uniqueProfile(ps)
	want := 2
	if got := len(ups); got != want {
		t.Errorf("UniqueProfile size = %v, want %v", got, want)
	}
}

func TestConfig_AddProfile(t *testing.T) {
	c := Load()
	c.AddProfile(APoint{
		Endpoint: "http://127.0.0.1:8082",
		Key:      "nat-key-1",
		Token:    "nat-token-1",
	})
	c.AddProfile(APoint{
		Endpoint: "http://127.0.0.1:8082",
		Key:      "nat-key-1",
	})
	c.AddProfile(APoint{
		Endpoint: "http://127.0.0.1:8082",
		Key:      "nat-key-1",
		Token:    "nat-token-2",
	})
	c.ClearActiveToken()
	want := 2
	if got := len(c.Scopes["http://127*0*0*1:8082"].Profiles); got != want {
		t.Errorf("Profile size = %v, want %v", got, want)
	}
}
