package config

import (
	"reflect"
	"testing"
)

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

func Test_compareStrs(t *testing.T) {
	type args struct {
		a string
		b string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"exact", args{"d65c62a114b710fe", "d65c62a114b710fe"}, 16},
		{"different", args{"d65c62a114b710fe", "nat"}, 0},
		{"partial", args{"d65c62a114b710fe", "d65c"}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareStrs(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("compareStrs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_min(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"simple", args{12, 3}, 3},
		{"equal", args{234, 234}, 234},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := min(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("min() = %v, want %v", got, tt.want)
			}
		})
	}
}
