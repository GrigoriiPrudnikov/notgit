package utils

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedParams []string
		expectedOpts   map[string]string
	}{
		{
			name:           "no args",
			args:           []string{},
			expectedParams: []string{},
			expectedOpts:   map[string]string{},
		},
		{
			name:           "one option with equal sign",
			args:           []string{"-opt=value"},
			expectedParams: []string{},
			expectedOpts:   map[string]string{"opt": "value"},
		},
		{
			name:           "one option with equal sign and 2 dashes",
			args:           []string{"--opt=value"},
			expectedParams: []string{},
			expectedOpts:   map[string]string{"opt": "value"},
		},
		{
			name:           "positional params only",
			args:           []string{"file1", "file2"},
			expectedParams: []string{"file1", "file2"},
			expectedOpts:   map[string]string{},
		},
		{
			name:           "mixed params and opts",
			args:           []string{"file.txt", "-a=1", "other", "--mode=fast"},
			expectedParams: []string{"file.txt", "other"},
			expectedOpts:   map[string]string{"a": "1", "mode": "fast"},
		},
		{
			name:           "option without value",
			args:           []string{"--flag"},
			expectedParams: []string{},
			expectedOpts:   map[string]string{"flag": "true"},
		},
		{
			name:           "multiple short opts",
			args:           []string{"-a=1", "-b=2", "-c=3"},
			expectedParams: []string{},
			expectedOpts:   map[string]string{"a": "1", "b": "2", "c": "3"},
		},
		{
			name:           "duplicate opts override",
			args:           []string{"--x=1", "--x=2"},
			expectedParams: []string{},
			expectedOpts:   map[string]string{"x": "2"},
		},
		{
			name:           "combined input",
			args:           []string{"some_param", "-a=1", "-b", "2", "--option_c", "3", "--option_d=4"},
			expectedParams: []string{"some_param"},
			expectedOpts:   map[string]string{"a": "1", "b": "2", "option_c": "3", "option_d": "4"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := ParseArgs(tt.args)
			if err != nil {
				t.Errorf("ParseArgs() error = %v", err)
				return
			}

			if !isEqualSlice(parsed.Params, tt.expectedParams) {
				t.Errorf("ParseArgs() params = %v, want %v", parsed.Params, tt.expectedParams)
			}

			if !isEqualMap(parsed.Opts, tt.expectedOpts) {
				t.Errorf("ParseArgs() opts = %v, want %v", parsed.Opts, tt.expectedOpts)
			}
		})
	}
}
