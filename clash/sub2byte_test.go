package clash

import (
	"net/url"
	"strings"
	"testing"
)

func TestHttpGet(t *testing.T) {
	type args struct {
		u string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				u: `https://qcrane-1.xyz/api/v1/client/subscribe?token=&flag=v2rayng`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt.args.u = strings.TrimSpace(tt.args.u)
		if _, err := url.Parse(tt.args.u); err != nil {
			t.Errorf("parse err in url %q, %s", tt.args.u, err)
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := HttpGet(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("HttpGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("httpget() got = %v", got)
			p := ParseContent(got)
			t.Logf("parse content found %d proxies", len(p))
		})
	}
}
