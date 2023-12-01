package set

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestValidateLogLevelArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{name: "valid args", args: []string{"osd", "10"}, wantErr: false},
		{name: "invalid ceph subsystem", args: []string{"test", "10"}, wantErr: true},
		{name: "invalid log level", args: []string{"osd", "101"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateLogLevelArgs(&cobra.Command{}, tt.args); (err != nil) != tt.wantErr {
				t.Errorf("validateLogLevelArgs error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
