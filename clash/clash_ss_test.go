package clash

import (
	"fmt"
	"testing"
)

func Test_buildSS(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want ClashSS
	}{
		{
			name: "test kycloud",
			args: args{
				s: "ss://#(SS)上海CU->美国7 GIA CN2 Netflix hulu HBO",
			},
		},
		{
			name: "test kycloud 2",
			args: args{
				s: "ss://#(SS)上海CN2->日本15 Sakura",
			},
		},
		{
			name: "test qcrane",
			args: args{
				s: "ss://@tunnel.nodelinks.xyz:30201#%5BSS%5D%20%E9%A6%99%E6%B8%AF%E9%9B%86%E7%BE%A4%288%E8%8A%82%E7%82%B9%29",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildSS(tt.args.s)
			fmt.Println(got)
		})
	}
}
