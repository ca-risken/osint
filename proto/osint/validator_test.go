package osint

import (
	"testing"
)

//Osint DataSource

func TestValidate_ListOsintRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListOsintRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListOsintRequest{ProjectId: 1001},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &ListOsintRequest{},
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

func TestValidate_GetOsintRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetOsintRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetOsintRequest{ProjectId: 1001, OsintId: 1002},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &GetOsintRequest{OsintId: 1002},
			wantErr: true,
		},
		{
			name:    "NG required(osint_data_source_id)",
			input:   &GetOsintRequest{ProjectId: 1001},
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

func TestValidate_PutOsintRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutOsintRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutOsintRequest{ProjectId: 1001, Osint: &OsintForUpsert{Name: "name", ProjectId: 1001}},
			wantErr: false,
		},
		{
			name:    "NG Required(Osint)",
			input:   &PutOsintRequest{ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Required(ProjectId)",
			input:   &PutOsintRequest{Osint: &OsintForUpsert{Name: "name", ProjectId: 1001}},
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

func TestValidate_DeleteOsintRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteOsintRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteOsintRequest{ProjectId: 1001, OsintId: 1002},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteOsintRequest{OsintId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Required(osint_data_source_id)",
			input:   &DeleteOsintRequest{ProjectId: 1001},
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

//OsintDataSource DataSource

func TestValidate_ListOsintDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListOsintDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListOsintDataSourceRequest{ProjectId: 1001},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &ListOsintDataSourceRequest{},
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

func TestValidate_GetOsintDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetOsintDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetOsintDataSourceRequest{ProjectId: 1001, OsintDataSourceId: 1002},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &GetOsintDataSourceRequest{OsintDataSourceId: 1002},
			wantErr: true,
		},
		{
			name:    "NG required(osint_data_source_id)",
			input:   &GetOsintDataSourceRequest{ProjectId: 1001},
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

func TestValidate_PutOsintDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutOsintDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutOsintDataSourceRequest{ProjectId: 1001, OsintDataSource: &OsintDataSourceForUpsert{Name: "name", Description: "description", MaxScore: 10.0}},
			wantErr: false,
		},
		{
			name:    "NG Required(OsintDataSource)",
			input:   &PutOsintDataSourceRequest{ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Required(ProjectId)",
			input:   &PutOsintDataSourceRequest{OsintDataSource: &OsintDataSourceForUpsert{Name: "name", Description: "description", MaxScore: 10.0}},
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

func TestValidate_DeleteOsintDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteOsintDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteOsintDataSourceRequest{ProjectId: 1001, OsintDataSourceId: 1002},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteOsintDataSourceRequest{OsintDataSourceId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Required(osint_data_source_id)",
			input:   &DeleteOsintDataSourceRequest{ProjectId: 1001},
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

//RelOsintDataSource DataSource

func TestValidate_ListRelOsintDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListRelOsintDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListRelOsintDataSourceRequest{ProjectId: 1001},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &ListRelOsintDataSourceRequest{},
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

func TestValidate_GetRelOsintDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetRelOsintDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetRelOsintDataSourceRequest{ProjectId: 1001, RelOsintDataSourceId: 1002},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &GetRelOsintDataSourceRequest{RelOsintDataSourceId: 1002},
			wantErr: true,
		},
		{
			name:    "NG required(rel_osint_data_source_id)",
			input:   &GetRelOsintDataSourceRequest{ProjectId: 1001},
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

func TestValidate_PutRelOsintDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutRelOsintDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutRelOsintDataSourceRequest{ProjectId: 1001, RelOsintDataSource: &RelOsintDataSourceForUpsert{ProjectId: 1001, OsintId: 1, OsintDataSourceId: 1, ResourceType: "domain", ResourceName: "hoge_name"}},
			wantErr: false,
		},
		{
			name:    "NG Required(RelOsintDataSource)",
			input:   &PutRelOsintDataSourceRequest{ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != rel_osint_data_source.project_id)",
			input:   &PutRelOsintDataSourceRequest{ProjectId: 1002, RelOsintDataSource: &RelOsintDataSourceForUpsert{ProjectId: 1001, OsintId: 1, OsintDataSourceId: 1, ResourceType: "domain", ResourceName: "hoge_name"}},
			wantErr: true,
		},
		{
			name:    "NG Required(ProjectId)",
			input:   &PutRelOsintDataSourceRequest{RelOsintDataSource: &RelOsintDataSourceForUpsert{ProjectId: 1001, OsintId: 1, OsintDataSourceId: 1, ResourceType: "domain", ResourceName: "hoge_name"}},
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

func TestValidate_DeleteRelOsintDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteRelOsintDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteRelOsintDataSourceRequest{ProjectId: 1001, RelOsintDataSourceId: 1002},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteRelOsintDataSourceRequest{RelOsintDataSourceId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Required(osint_data_source_id)",
			input:   &DeleteRelOsintDataSourceRequest{ProjectId: 1001},
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

func TestValidate_StartOsintRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *StartOsintRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &StartOsintRequest{ProjectId: 1001, RelOsintDataSourceId: 1002},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &StartOsintRequest{RelOsintDataSourceId: 1002},
			wantErr: true,
		},
		{
			name:    "NG Required(rel_osint_data_source_id)",
			input:   &StartOsintRequest{ProjectId: 1001},
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

func TestValidate_OsintForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *OsintForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &OsintForUpsert{Name: "name", ProjectId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(name)",
			input:   &OsintForUpsert{ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Length(name)",
			input:   &OsintForUpsert{Name: "123456789012345678901234567890123456789012345678901", ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Required(ProjectId)",
			input:   &OsintForUpsert{Name: "name"},
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

func TestValidate_OsintDataSourceForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *OsintDataSourceForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &OsintDataSourceForUpsert{Name: "name", Description: "description", MaxScore: 100},
			wantErr: false,
		},
		{
			name:    "NG Length(name)",
			input:   &OsintDataSourceForUpsert{Name: "123456789012345678901234567890123456789012345678901", Description: "description", MaxScore: 100},
			wantErr: true,
		},
		{
			name:    "NG Required(name)",
			input:   &OsintDataSourceForUpsert{Description: "description", MaxScore: 100},
			wantErr: true,
		},
		{
			name:    "NG Length(description)",
			input:   &OsintDataSourceForUpsert{Name: "name", Description: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012", MaxScore: 100},
			wantErr: true,
		},
		{
			name:    "NG Required(description)",
			input:   &OsintDataSourceForUpsert{Name: "name", MaxScore: 100},
			wantErr: true,
		},
		{
			name:    "NG Num Over(Max Score)",
			input:   &OsintDataSourceForUpsert{Name: "name", Description: "description", MaxScore: 100000},
			wantErr: true,
		},
		{
			name:    "NG Num Under(Max Score)",
			input:   &OsintDataSourceForUpsert{Name: "name", Description: "description", MaxScore: -1.0},
			wantErr: true,
		},
		{
			name:    "NG Required(Max Score)",
			input:   &OsintDataSourceForUpsert{Name: "name", Description: "description"},
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

func TestValidate_RelOsintDataSourceForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *RelOsintDataSourceForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &RelOsintDataSourceForUpsert{ProjectId: 1001, OsintId: 1001, OsintDataSourceId: 1001, ResourceType: "hoge_type", ResourceName: "hoge_name"},
			wantErr: false,
		},
		{
			name:    "NG Required(osint_id)",
			input:   &RelOsintDataSourceForUpsert{ProjectId: 1001, OsintDataSourceId: 1001, ResourceType: "hoge_type", ResourceName: "hoge_name"},
			wantErr: true,
		},
		{
			name:    "NG Required(osint_data_source_id)",
			input:   &RelOsintDataSourceForUpsert{ProjectId: 1001, OsintId: 1001, ResourceType: "hoge_type", ResourceName: "hoge_name"},
			wantErr: true,
		},
		{
			name:    "NG Required(project_id)",
			input:   &RelOsintDataSourceForUpsert{OsintId: 1001, OsintDataSourceId: 1001, ResourceType: "hoge_type", ResourceName: "hoge_name"},
			wantErr: true,
		},
		{
			name:    "Too long(resource_type)",
			input:   &RelOsintDataSourceForUpsert{ProjectId: 1001, OsintId: 1001, OsintDataSourceId: 1001, ResourceType: "123456789012345678901234567890123456789012345678901", ResourceName: "hoge_name"},
			wantErr: true,
		},
		{
			name:    "NG Required(resource_type)",
			input:   &RelOsintDataSourceForUpsert{ProjectId: 1001, OsintId: 1001, OsintDataSourceId: 1001, ResourceName: "hoge_name"},
			wantErr: true,
		},
		{
			name:    "Too long(resource_name)",
			input:   &RelOsintDataSourceForUpsert{ProjectId: 1001, OsintId: 1001, OsintDataSourceId: 1001, ResourceType: "hoge_type", ResourceName: "123456789012345678901234567890123456789012345678901123456789012345678901234567890123456789012345678901123456789012345678901234567890123456789012345678901123456789012345678901234567890123456789012345678901"},
			wantErr: true,
		},
		{
			name:    "NG Too small scan_at",
			input:   &RelOsintDataSourceForUpsert{ProjectId: 1001, OsintId: 1001, OsintDataSourceId: 1001, ResourceType: "hoge_type", ResourceName: "hoge_name", ScanAt: -1},
			wantErr: true,
		},
		{
			name:    "NG Too large scan_at",
			input:   &RelOsintDataSourceForUpsert{ProjectId: 1001, OsintId: 1001, OsintDataSourceId: 1001, ResourceType: "hoge_type", ResourceName: "hoge_name", ScanAt: 253402268400},
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
