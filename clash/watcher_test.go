package clash

import (
	"os"
	"testing"
)

func TestUpdate(t *testing.T) {
	os.Setenv("SUB_BLACKLIST", "")
	os.Setenv("SUB_WHITELIST", "回国,海外")
	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Update("")
		})
	}
}
