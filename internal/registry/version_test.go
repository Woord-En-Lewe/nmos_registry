package registry

import "testing"

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Version
		wantErr bool
	}{
		{
			name:  "valid version",
			input: "1441700172:318426300",
			want:  &Version{Seconds: 1441700172, Nanos: 318426300},
		},
		{
			name:  "zero version",
			input: "0:0",
			want:  &Version{Seconds: 0, Nanos: 0},
		},
		{
			name:    "invalid format - no colon",
			input:   "1441700172.318426300",
			wantErr: true,
		},
		{
			name:    "invalid format - too many parts",
			input:   "1441700172:318426300:1",
			wantErr: true,
		},
		{
			name:    "invalid seconds",
			input:   "abc:318426300",
			wantErr: true,
		},
		{
			name:    "invalid nanos",
			input:   "1441700172:abc",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVersion(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Seconds != tt.want.Seconds {
				t.Errorf("ParseVersion(%q).Seconds = %v, want %v", tt.input, got.Seconds, tt.want.Seconds)
			}
			if !tt.wantErr && got.Nanos != tt.want.Nanos {
				t.Errorf("ParseVersion(%q).Nanos = %v, want %v", tt.input, got.Nanos, tt.want.Nanos)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name    string
		v1      string
		v2      string
		want    int
		wantErr bool
	}{
		{
			name: "v1 later than v2 (seconds)",
			v1:   "1441700173:0",
			v2:   "1441700172:0",
			want: 1,
		},
		{
			name: "v1 earlier than v2 (seconds)",
			v1:   "1441700172:0",
			v2:   "1441700173:0",
			want: -1,
		},
		{
			name: "v1 later than v2 (nanos)",
			v1:   "1441700172:500000000",
			v2:   "1441700172:318426300",
			want: 1,
		},
		{
			name: "v1 earlier than v2 (nanos)",
			v1:   "1441700172:100000000",
			v2:   "1441700172:318426300",
			want: -1,
		},
		{
			name: "equal versions",
			v1:   "1441700172:318426300",
			v2:   "1441700172:318426300",
			want: 0,
		},
		{
			name:    "invalid v1",
			v1:      "invalid",
			v2:      "1441700172:318426300",
			wantErr: true,
		},
		{
			name:    "invalid v2",
			v1:      "1441700172:318426300",
			v2:      "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompareVersions(tt.v1, tt.v2)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareVersions(%q, %q) error = %v, wantErr %v", tt.v1, tt.v2, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CompareVersions(%q, %q) = %v, want %v", tt.v1, tt.v2, got, tt.want)
			}
		})
	}
}

func TestIsVersionEarlier(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want bool
	}{
		{
			name: "v1 is earlier",
			v1:   "1441700172:0",
			v2:   "1441700173:0",
			want: true,
		},
		{
			name: "v1 is later",
			v1:   "1441700173:0",
			v2:   "1441700172:0",
			want: false,
		},
		{
			name: "equal",
			v1:   "1441700172:318426300",
			v2:   "1441700172:318426300",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsVersionEarlier(tt.v1, tt.v2)
			if err != nil {
				t.Errorf("IsVersionEarlier error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("IsVersionEarlier(%q, %q) = %v, want %v", tt.v1, tt.v2, got, tt.want)
			}
		})
	}
}