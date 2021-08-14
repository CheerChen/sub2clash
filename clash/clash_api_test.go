package clash

import (
	"fmt"
	"testing"
)

func TestGetDelayMap(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test getdelaymap",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetProxies(); (err != nil) != tt.wantErr {
				t.Errorf("GetProxies() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(regionList)
		})
	}
}
