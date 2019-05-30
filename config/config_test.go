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
		{"default active", fields{DefaultConfig().Scopes, DefaultConfig().Active}, APoint{DefaultEndpoint, "", "", DefaultAppKey, ""}},
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

func TestConfig_GetProfile(t *testing.T) {
	type fields struct {
		Scopes map[string]Scope
		Active string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Profile
		wantErr bool
	}{
		{"partial match", fields{DefaultConfig().Scopes, DefaultConfig().Active}, args{"def"}, &Profile{"default", "", "", DefaultAppKey, ""}, false},
		{"not found", fields{DefaultConfig().Scopes, DefaultConfig().Active}, args{"nat"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Scopes: tt.fields.Scopes,
				Active: tt.fields.Active,
			}
			got, err := c.GetProfile(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.GetProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.GetProfile() = %v, want %v", got, tt.want)
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
      email: ""
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
		{"empty", fields{"nat-endpoint", []Profile{}}, args{Profile{"natid", "Nat", "nat@nat.com", "nat-key", "nat-token"}}, Profile{"natid", "Nat", "nat@nat.com", "nat-key", "nat-token"}},
		{"exists", fields{"nat-endpoint", []Profile{{"aid", "Nat", "nat@nat.com", "nat-key", "nat-token"}}}, args{Profile{"natid", "Nat", "nat@nat.com", "nat-key", "nat-token"}}, Profile{"aid", "Nat", "nat@nat.com", "nat-key", "nat-token"}},
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
		Email string
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
		{"equal", fields{"anyid", "Nat1", "nat@nat.com", "nat-key", "nat-token"}, args{Profile{"anoterid", "Nat2", "nat2@nat.com", "nat-key", "nat-token"}}, true},
		{"different key", fields{"anyid", "Nat", "nat@nat.com", "nat-key", "nat-token"}, args{Profile{"anoterid", "Nat", "nat@nat.com", "a-key", "nat-token"}}, false},
		{"different token", fields{"anyid", "Nat", "nat@nat.com", "nat-key", "nat-token"}, args{Profile{"anyid", "Nat", "nat@nat.com", "nat-key", "a-token"}}, false},
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
