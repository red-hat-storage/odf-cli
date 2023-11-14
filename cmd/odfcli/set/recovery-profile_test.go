package set

import "testing"

func Test_isValidArg(t *testing.T) {
	type args struct {
		args string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Valid Input with 'high_client_ops' as arg", args{args: "high_client_ops"}, true},
		{"Valid Input with 'balanced' as arg", args{args: "balanced"}, true},
		{"Invalid input for mclock profile", args{args: "balance"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidArg(tt.args.args); got != tt.want {
				t.Errorf("isValidArg() = %v, want %v", got, tt.want)
			}
		})
	}
}
