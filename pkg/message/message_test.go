package message

import (
	"reflect"
	"testing"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name    string
		input   *OsintQueueMessage
		wantErr bool
	}{
		{
			name:  "OK (subdomain)",
			input: &OsintQueueMessage{OsintID: 1001, OsintDataSourceID: 1001, RelOsintDataSourceID: 1001, DataSource: "osint:subdomain", ProjectID: 1001, ResourceName: "hogefuga", ResourceType: "domain"},
		},
		{
			name:  "OK (website)",
			input: &OsintQueueMessage{OsintID: 1001, OsintDataSourceID: 1001, RelOsintDataSourceID: 1001, DataSource: "osint:website", ProjectID: 1001, ResourceName: "hogefuga", ResourceType: "domain"},
		},
		{
			name:    "NG Required(OsintID)",
			input:   &OsintQueueMessage{OsintDataSourceID: 1001, RelOsintDataSourceID: 1001, DataSource: "osint:website", ProjectID: 1001, ResourceName: "hogefuga", ResourceType: "domain"},
			wantErr: true,
		},
		{
			name:    "NG Required(OsintDataSourceID)",
			input:   &OsintQueueMessage{OsintID: 1001, RelOsintDataSourceID: 1001, DataSource: "osint:website", ProjectID: 1001, ResourceName: "hogefuga", ResourceType: "domain"},
			wantErr: true,
		},
		{
			name:    "NG Required(RelOsintDataSourceID)",
			input:   &OsintQueueMessage{OsintID: 1001, OsintDataSourceID: 1001, DataSource: "osint:website", ProjectID: 1001, ResourceName: "hogefuga", ResourceType: "domain"},
			wantErr: true,
		},
		{
			name:    "NG Required(DataSource)",
			input:   &OsintQueueMessage{OsintID: 1001, OsintDataSourceID: 1001, RelOsintDataSourceID: 1001, ProjectID: 1001, ResourceName: "hogefuga", ResourceType: "domain"},
			wantErr: true,
		},
		{
			name:    "NG Unknown(DataSource)",
			input:   &OsintQueueMessage{OsintID: 1001, OsintDataSourceID: 1001, RelOsintDataSourceID: 1001, DataSource: "osint:unknown", ProjectID: 1001, ResourceName: "hogefuga", ResourceType: "domain"},
			wantErr: true,
		},
		{
			name:    "NG Required(ProjectID)",
			input:   &OsintQueueMessage{OsintID: 1001, OsintDataSourceID: 1001, RelOsintDataSourceID: 1001, DataSource: "osint:website", ResourceName: "hogefuga", ResourceType: "domain"},
			wantErr: true,
		},
		{
			name:    "NG Required(ResourceName)",
			input:   &OsintQueueMessage{OsintID: 1001, OsintDataSourceID: 1001, RelOsintDataSourceID: 1001, DataSource: "osint:website", ProjectID: 1001, ResourceType: "domain"},
			wantErr: true,
		},
		{
			name:    "NG Required(ResourceType)",
			input:   &OsintQueueMessage{OsintID: 1001, OsintDataSourceID: 1001, RelOsintDataSourceID: 1001, DataSource: "osint:website", ProjectID: 1001, ResourceName: "hogefuga"},
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

func TestParseMessage(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    *OsintQueueMessage
		wantErr bool
	}{
		{
			name:  "OK",
			input: `{"osint_id":1001, "osint_data_source_id":1001, "rel_osint_data_source_id":1001, "data_source":"osint:subdomain", "project_id":1001, "resource_name":"hogefuga", "resource_type":"domain"}`,
			want:  &OsintQueueMessage{OsintID: 1001, OsintDataSourceID: 1001, RelOsintDataSourceID: 1001, DataSource: "osint:subdomain", ProjectID: 1001, ResourceName: "hogefuga", ResourceType: "domain"},
		},
		{
			name:  "OK(exist scan_only)",
			input: `{"osint_id":1001, "osint_data_source_id":1001, "rel_osint_data_source_id":1001, "data_source":"osint:subdomain", "project_id":1001, "resource_name":"hogefuga", "resource_type":"domain", "scan_only":"true"}`,
			want:  &OsintQueueMessage{OsintID: 1001, OsintDataSourceID: 1001, RelOsintDataSourceID: 1001, DataSource: "osint:subdomain", ProjectID: 1001, ResourceName: "hogefuga", ResourceType: "domain", ScanOnly: true},
		},
		{
			name:    "NG Json parse erroro",
			input:   `aaaaa`,
			wantErr: true,
		},
		{
			name:    "NG Invalid message(required parammeter)",
			input:   `{}`,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ParseMessage(c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error occured, wantErr=%t, err=%+v", c.wantErr, err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpaeted response, want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
