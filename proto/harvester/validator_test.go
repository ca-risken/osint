package harvester

import (
	"testing"
)

func TestValidate_InvokeScanRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *InvokeScanRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &InvokeScanRequest{ResourceName: "hogehoge.com", ResourceType: "Domain"},
			wantErr: false,
		},
		{
			name:    "NG Required(resource_name)",
			input:   &InvokeScanRequest{ResourceType: "Domain"},
			wantErr: true,
		},
		{
			name:    "NG Required(resource_type)",
			input:   &InvokeScanRequest{ResourceName: "hogehoge.com"},
			wantErr: true,
		},
		{
			name:    "NG Length(resource_name)",
			input:   &InvokeScanRequest{ResourceName: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901", ResourceType: "Domain"},
			wantErr: true,
		},
		{
			name:    "NG Length(resource_type)",
			input:   &InvokeScanRequest{ResourceName: "hogehoge.com", ResourceType: "123456789012345678901234567890123456789012345678901"},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}
