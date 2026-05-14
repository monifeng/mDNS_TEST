package netutil

import (
	"testing"
)

func TestParsePorts(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
		wantErr  bool
	}{
		{
			name:     "单个端口",
			input:    "80",
			expected: []int{80},
			wantErr:  false,
		},
		{
			name:     "多个端口逗号分隔",
			input:    "80,443,8080",
			expected: []int{80, 443, 8080},
			wantErr:  false,
		},
		{
			name:     "端口范围",
			input:    "10-20",
			expected: []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			wantErr:  false,
		},
		{
			name:     "混合范围和单端口",
			input:    "9,445,548,5000",
			expected: []int{9, 445, 548, 5000},
			wantErr:  false,
		},
		{
			name:     "无效端口0",
			input:    "0",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "无效端口65536",
			input:    "65536",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "无效范围倒序",
			input:    "100-10",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePorts(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePorts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				for _, p := range tt.expected {
					if !got[p] {
						t.Errorf("ParsePorts() missing expected port %d", p)
					}
				}
			}
		})
	}
}

func TestIPInCIDR(t *testing.T) {
	ipnet, err := ParseCIDR("192.168.1.0/24")
	if err != nil {
		t.Fatalf("ParseCIDR failed: %v", err)
	}
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{
			name: "网段内",
			ip:   "192.168.1.10",
			want: true,
		},
		{
			name: "网段外",
			ip:   "192.168.2.10",
			want: false,
		},
		{
			name: "无效IP",
			ip:   "999.999.999.999",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IPInCIDR(tt.ip, ipnet); got != tt.want {
				t.Errorf("IPInCIDR() = %v, want %v", got, tt.want)
			}
		})
	}
}
