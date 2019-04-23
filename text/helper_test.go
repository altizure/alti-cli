package text

import (
	"testing"
)

func TestContains(t *testing.T) {
	type args struct {
		a []string
		s string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
		{"empty", args{[]string{}, ""}, -1, false},
		{"found", args{[]string{"abc", "nat", "pi"}, "nat"}, 1, true},
		{"not found", args{[]string{"abc", "nat", "pi"}, "bunny"}, -1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Contains(tt.args.a, tt.args.s)
			if got != tt.want {
				t.Errorf("Contains() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Contains() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestBestMatch(t *testing.T) {
	type args struct {
		a   []string
		s   string
		def string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{[]string{}, "", ""}, ""},
		{"exact match", args{[]string{"abc", "nat", "pi", "nata"}, "nat", ""}, "nat"},
		{"regex match", args{[]string{"abc", "nat", "pi", "nata"}, "na", ""}, "nat"},
		{"regex match", args{[]string{"abc", "nat", "pi", "nata"}, "ata", ""}, "nata"},
		{"regex match", args{[]string{"abc", "nat", "pi", "nata", "natal"}, "ta", ""}, "nata"},
		{"case insensitive", args{[]string{"abc", "Nat", "Pi", "UST", "Jan"}, "us", ""}, "UST"},
		{"special regex characters", args{[]string{"abc", "nat", "pi", "nata", "natal"}, "()|[]{", ""}, ""},
		{"no match", args{[]string{"abc", "nat", "pi", "nata"}, "123", "bunny"}, "bunny"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BestMatch(tt.args.a, tt.args.s, tt.args.def); got != tt.want {
				t.Errorf("BestMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
