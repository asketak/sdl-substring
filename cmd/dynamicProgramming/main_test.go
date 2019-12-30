package main

import (
	"strings"
	"testing"
)

func TestLCSubstring(t *testing.T) {
	type args struct {
		s string
		t string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"samestring", args{"hahoj", "hahoj"},"hahoj"},
		{"nocommon", args{"abcd", "efgh"},""},

		{"emptyfirst", args{"hahoj", ""},""},
		{"emptysecond", args{"", "hahoj"},""},

		{"prefix", args{"aaaa", "aahj"},"aa"},
		{"prefix2", args{"aahj", "aaaa"},"aa"},
		{"prefix3", args{"aa", "aasa"},"aa"},
		{"prefix4", args{"aagg", "aa"},"aa"},

		{"sufix", args{"hjaa", "aaaa"},"aa"},
		{"sufix2", args{"aaaa", "hjaa"},"aa"},
		{"sufix3", args{"asaa", "aa"},"aa"},
		{"sufix4", args{"aa", "ggaa"},"aa"},

		{"middle", args{"aaa", "ggaaahh"},"aaa"},
		{"middle2", args{"ggaaahh", "aaa"},"aaa"},
		{"middle3", args{"aaa", "ggaaaahh"},"aaa"},

		{"both ends", args{"aaa", "aaaJJaaa"},"aaa"},

		{"mutliple times", args{"gagagaga", "gaga"},"gaga"},
		{"three substrings same length", args{"hahaIgagaLhehe", "heheMgagaOhaha"},"haha"},

		{"small match than long", args{"hahaIgagaLha", "fffffhahafffff"},"haha"}, // finds first biggest substring
		{"long match than small", args{"haIgagaLhaha", "fffffhahafffff"},"haha"}, // finds first biggest substring

		{"no match", args{"hhhhhh", "faaaff"},""},
		{"single char match ", args{"hhahhh", "faaaff"},"a"},
		{"hundred char match ", args{strings.Repeat("na", 100), strings.Repeat("ha", 100)},"a"},
//		{"thousands char match ", args{strings.Repeat("na", 1000), strings.Repeat("ha", 1000)},"a"},
//		{"thousands char match ", args{strings.Repeat("na", 10000), strings.Repeat("ha", 10000)},"a"},
//		{"hundred thousand char match ", args{strings.Repeat("na", 100000), strings.Repeat("ha", 100000)},"a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LCSubstring(tt.args.s, tt.args.t); got != tt.want {
				t.Errorf(" test: %v: LCSubstring() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
