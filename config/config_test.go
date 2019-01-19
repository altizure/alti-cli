package config

import (
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		{"default config", DefaultConfig()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Load(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetActive(t *testing.T) {
	type fields struct {
		Scopes map[string]Scope
		Active string
	}
	tests := []struct {
		name   string
		fields fields
		want   APoint
	}{
		{"default active", fields{DefaultConfig().Scopes, DefaultConfig().Active}, APoint{DefaultEndpoint, "", DefaultAppKey, ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Scopes: tt.fields.Scopes,
				Active: tt.fields.Active,
			}
			if got := c.GetActive(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.GetActive() = %v, want %v", got, tt.want)
			}
		})
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
	c.ClearActiveToken(false)
	want := 2
	if got := len(c.Scopes["http://127*0*0*1:8082"].Profiles); got != want {
		t.Errorf("Profile size = %v, want %v", got, want)
	}
}

func TestConfig_String(t *testing.T) {
	type fields struct {
		Scopes map[string]Scope
		Active string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"default string", fields{DefaultConfig().Scopes, DefaultConfig().Active}, `scopes:
  https://api*altizure*com:
    endpoint: https://api.altizure.com
    profiles:
    - id: default
      name: ""
      key: Ah8bOakrkmSl2FA9OCbT8EnFOUrPwOOZ7HQxZm6
      token: ""
active: default
`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Scopes: tt.fields.Scopes,
				Active: tt.fields.Active,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("Config.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScope_Add(t *testing.T) {
	type fields struct {
		Endpoint string
		Profiles []Profile
	}
	type args struct {
		p Profile
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Profile
	}{
		{"empty", fields{"nat-endpoint", []Profile{}}, args{Profile{"natid", "Nat", "nat-key", "nat-token"}}, Profile{"natid", "Nat", "nat-key", "nat-token"}},
		{"exists", fields{"nat-endpoint", []Profile{{"aid", "Nat", "nat-key", "nat-token"}}}, args{Profile{"natid", "Nat", "nat-key", "nat-token"}}, Profile{"aid", "Nat", "nat-key", "nat-token"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scope{
				Endpoint: tt.fields.Endpoint,
				Profiles: tt.fields.Profiles,
			}
			if got := s.Add(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scope.Add() = %v, want %v", got, tt.want)
			}
			if got := len(s.Profiles); got != 1 {
				t.Errorf("Scope.Add() = %v, want %v", got, 1)
			}
		})
	}
}

func TestProfile_Equal(t *testing.T) {
	type fields struct {
		ID    string
		Name  string
		Key   string
		Token string
	}
	type args struct {
		o Profile
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"equal", fields{"anyid", "Nat1", "nat-key", "nat-token"}, args{Profile{"anoterid", "Nat2", "nat-key", "nat-token"}}, true},
		{"different key", fields{"anyid", "Nat", "nat-key", "nat-token"}, args{Profile{"anoterid", "Nat", "a-key", "nat-token"}}, false},
		{"different token", fields{"anyid", "Nat", "nat-key", "nat-token"}, args{Profile{"anyid", "Nat", "nat-key", "a-token"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Profile{
				ID:    tt.fields.ID,
				Key:   tt.fields.Key,
				Token: tt.fields.Token,
			}
			if got := p.Equal(tt.args.o); got != tt.want {
				t.Errorf("Profile.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_endpointToKey(t *testing.T) {
	type args struct {
		ep string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"normal", args{"https://api.altizure.com/graphql"}, "https://api*altizure*com"},
		{"empty", args{""}, "https://api*altizure*com"},
		{"invalid", args{"nat"}, "https://api*altizure*com"},
		{"no scheme", args{"nat.com"}, "https://api*altizure*com"},
		{"with port", args{"http://nat.com:12345/abcde"}, "http://nat*com:12345"},
		{"localhost", args{"http://127.0.0.1:8082"}, "http://127*0*0*1:8082"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := endpointToKey(tt.args.ep); got != tt.want {
				t.Errorf("endpointToKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uniqueProfile(t *testing.T) {
	type args struct {
		ps []Profile
	}
	tests := []struct {
		name string
		args args
		want []Profile
	}{
		{"simple", args{[]Profile{{"id1", "n1", "k1", "t1"}, {"id2", "n1", "k2", "t2"}, {"id1", "n1", "k1", "t1"}}}, []Profile{{"id1", "n1", "k1", "t1"}, {"id2", "n1", "k2", "t2"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := uniqueProfile(tt.args.ps); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uniqueProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}
