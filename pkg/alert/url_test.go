package alert

import "testing"

func TestParseAndValidateURL(t *testing.T) {
	tests := []struct {
		name             string
		url              string
		supportedSchemes []string
		wantErr          bool
	}{
		{
			name:             "Empty URL",
			url:              "",
			supportedSchemes: defaultSchemes,
			wantErr:          true,
		},
		{
			name:             "Invalid URL Format",
			url:              "::invalid-url",
			supportedSchemes: defaultSchemes,
			wantErr:          true,
		},
		{
			name:             "No Protocol Specified",
			url:              "//example.com",
			supportedSchemes: defaultSchemes,
			wantErr:          true,
		},
		{
			name:             "Unsupported Protocol",
			url:              "htp://xyz.testdomain.abcdef",
			supportedSchemes: defaultSchemes,
			wantErr:          true,
		},
		{
			name:             "No Host",
			url:              "https://",
			supportedSchemes: defaultSchemes,
			wantErr:          true,
		},
		{
			name:             "Valid URL",
			url:              "https://vxn.dev",
			supportedSchemes: defaultSchemes,
			wantErr:          false,
		},
		{
			name:             "Custom Supported Schemes with Valid URL",
			url:              "ftp://vxn.dev",
			supportedSchemes: []string{"ftp"},
			wantErr:          false,
		},
		{
			name:             "Custom Supported Schemes with Invalid URL",
			url:              "https://vxn.dev",
			supportedSchemes: []string{"ftp"},
			wantErr:          true,
		},
		{
			name:             "No Supported Schemes Provided (nil)",
			url:              "https://vxn.dev",
			supportedSchemes: nil,
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseAndValidateURL(tt.url, tt.supportedSchemes)

			if tt.wantErr && err == nil {
				t.Error("expected an error but got none")
			} else if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
