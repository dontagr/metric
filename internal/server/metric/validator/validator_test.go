package validator

import "testing"

func TestIsValidMType(t *testing.T) {
	type args struct {
		mType string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid Counter Type",
			args: args{mType: "counter"},
			want: true,
		},
		{
			name: "Valid Gauge Type",
			args: args{mType: "gauge"},
			want: true,
		},
		{
			name: "Invalid Type",
			args: args{mType: "invalid_type"},
			want: false,
		},
		{
			name: "Empty Type",
			args: args{mType: ""},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidMType(tt.args.mType); got != tt.want {
				t.Errorf("IsValidMType() = %v, want %v", got, tt.want)
			}
		})
	}
}
