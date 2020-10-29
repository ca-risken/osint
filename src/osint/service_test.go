package main

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/CyberAgent/mimosa-osint-go/pkg/model"
	"github.com/CyberAgent/mimosa-osint-go/proto/osint"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

func TestListOsint(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockOsintRepository{}
	svc := osintService{repository: &mockDB}
	cases := []struct {
		name         string
		input        *osint.ListOsintRequest
		want         *osint.ListOsintResponse
		mockResponce *[]model.Osint
		mockError    error
	}{
		{
			name:  "OK",
			input: &osint.ListOsintRequest{ProjectId: 1001},
			want: &osint.ListOsintResponse{Osint: []*osint.Osint{
				{OsintId: 1001, ProjectId: 1001, ResourceType: "test_type", ResourceName: "test_name", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{OsintId: 1002, ProjectId: 1001, ResourceType: "test_type", ResourceName: "test_name2", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponce: &[]model.Osint{
				{OsintID: 1001, ProjectID: 1001, ResourceType: "test_type", ResourceName: "test_name", CreatedAt: now, UpdatedAt: now},
				{OsintID: 1002, ProjectID: 1001, ResourceType: "test_type", ResourceName: "test_name2", CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK Record not found",
			input:     &osint.ListOsintRequest{ProjectId: 1001},
			want:      &osint.ListOsintResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListOsint").Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListOsint(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetOsint(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockOsintRepository{}
	svc := osintService{repository: &mockDB}
	cases := []struct {
		name         string
		input        *osint.GetOsintRequest
		want         *osint.GetOsintResponse
		mockResponce *model.Osint
		mockError    error
	}{
		{
			name:  "OK",
			input: &osint.GetOsintRequest{ProjectId: 1001, OsintId: 1001},
			want: &osint.GetOsintResponse{Osint: &osint.Osint{
				OsintId: 1001, ProjectId: 1001, ResourceType: "test_type", ResourceName: "test_name", CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
			}},
			mockResponce: &model.Osint{
				OsintID: 1001, ProjectID: 1001, ResourceType: "test_type", ResourceName: "test_name", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "OK Record not found",
			input:     &osint.GetOsintRequest{ProjectId: 1001, OsintId: 1001},
			want:      &osint.GetOsintResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetOsint").Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.GetOsint(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutOsint(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockOsintRepository{}
	svc := osintService{repository: &mockDB}
	cases := []struct {
		name        string
		input       *osint.PutOsintRequest
		want        *osint.PutOsintResponse
		wantErr     bool
		mockUpdResp *model.Osint
		mockUpdErr  error
	}{
		{
			name:        "OK Update",
			input:       &osint.PutOsintRequest{ProjectId: 1001, Osint: &osint.OsintForUpsert{ResourceType: "test_type", ResourceName: "test_name", ProjectId: 1001, OsintId: 1001}},
			want:        &osint.PutOsintResponse{Osint: &osint.Osint{OsintId: 1001, ResourceType: "test_type", ResourceName: "test_name", ProjectId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpdResp: &model.Osint{OsintID: 1001, ResourceType: "test_type", ResourceName: "test_name", ProjectID: 1001, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Insert",
			input:       &osint.PutOsintRequest{ProjectId: 1001, Osint: &osint.OsintForUpsert{ResourceType: "test_type", ResourceName: "test_name", ProjectId: 1001}},
			want:        &osint.PutOsintResponse{Osint: &osint.Osint{OsintId: 1001, ResourceType: "test_type", ResourceName: "test_name", ProjectId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpdResp: &model.Osint{OsintID: 1001, ResourceType: "test_type", ResourceName: "test_name", ProjectID: 1001, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid parameter(required project_id)",
			input:   &osint.PutOsintRequest{Osint: &osint.OsintForUpsert{ResourceType: "test_type", ResourceName: "test_name", ProjectId: 1001, OsintId: 1001}},
			wantErr: true,
		},
		{
			name:        "NG DB error(UpsertOsint)",
			input:       &osint.PutOsintRequest{Osint: &osint.OsintForUpsert{ResourceType: "test_type", ResourceName: "test_name", ProjectId: 1001, OsintId: 1001}},
			wantErr:     true,
			mockUpdResp: nil,
			mockUpdErr:  errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockUpdResp != nil || c.mockUpdErr != nil {
				mockDB.On("UpsertOsint").Return(c.mockUpdResp, c.mockUpdErr).Once()
			}
			got, err := svc.PutOsint(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteOsint(t *testing.T) {
	var ctx context.Context
	mockDB := mockOsintRepository{}
	svc := osintService{repository: &mockDB}
	cases := []struct {
		name     string
		input    *osint.DeleteOsintRequest
		wantErr  bool
		mockResp error
	}{
		{
			name:     "OK",
			input:    &osint.DeleteOsintRequest{ProjectId: 1001, OsintId: 1001},
			wantErr:  false,
			mockResp: nil,
		},
		{
			name:     "NG DB error",
			input:    &osint.DeleteOsintRequest{ProjectId: 1001, OsintId: 1001},
			wantErr:  true,
			mockResp: errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("DeleteOsint").Return(c.mockResp).Once()
			_, err := svc.DeleteOsint(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
		})
	}
}

func TestListOsintDataSource(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockOsintRepository{}
	svc := osintService{repository: &mockDB}
	cases := []struct {
		name         string
		input        *osint.ListOsintDataSourceRequest
		want         *osint.ListOsintDataSourceResponse
		mockResponce *[]model.OsintDataSource
		mockError    error
	}{
		{
			name:  "OK",
			input: &osint.ListOsintDataSourceRequest{ProjectId: 1001, Name: "test"},
			want: &osint.ListOsintDataSourceResponse{OsintDataSource: []*osint.OsintDataSource{
				{OsintDataSourceId: 1001, Name: "test_name", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{OsintDataSourceId: 1002, Name: "test_name2", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponce: &[]model.OsintDataSource{
				{OsintDataSourceID: 1001, Name: "test_name", CreatedAt: now, UpdatedAt: now},
				{OsintDataSourceID: 1002, Name: "test_name2", CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK Record not found",
			input:     &osint.ListOsintDataSourceRequest{ProjectId: 1001, Name: "test"},
			want:      &osint.ListOsintDataSourceResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListOsintDataSource").Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListOsintDataSource(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetOsintDataSource(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockOsintRepository{}
	svc := osintService{repository: &mockDB}
	cases := []struct {
		name         string
		input        *osint.GetOsintDataSourceRequest
		want         *osint.GetOsintDataSourceResponse
		mockResponce *model.OsintDataSource
		mockError    error
	}{
		{
			name:  "OK",
			input: &osint.GetOsintDataSourceRequest{ProjectId: 1001, OsintDataSourceId: 1001},
			want: &osint.GetOsintDataSourceResponse{OsintDataSource: &osint.OsintDataSource{
				OsintDataSourceId: 1001, Name: "test_name", CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
			}},
			mockResponce: &model.OsintDataSource{
				OsintDataSourceID: 1001, Name: "test_name", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "OK Record not found",
			input:     &osint.GetOsintDataSourceRequest{ProjectId: 1001, OsintDataSourceId: 1001},
			want:      &osint.GetOsintDataSourceResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetOsintDataSource").Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.GetOsintDataSource(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutOsintDataSource(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockOsintRepository{}
	svc := osintService{repository: &mockDB}
	cases := []struct {
		name        string
		input       *osint.PutOsintDataSourceRequest
		want        *osint.PutOsintDataSourceResponse
		wantErr     bool
		mockUpdResp *model.OsintDataSource
		mockUpdErr  error
	}{
		{
			name:        "OK Update",
			input:       &osint.PutOsintDataSourceRequest{ProjectId: 1001, OsintDataSource: &osint.OsintDataSourceForUpsert{Name: "test_name", Description: "test_desc", MaxScore: 10.0, OsintDataSourceId: 1001}},
			want:        &osint.PutOsintDataSourceResponse{OsintDataSource: &osint.OsintDataSource{OsintDataSourceId: 1001, Name: "test_name", Description: "test_desc", MaxScore: 10.0, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpdResp: &model.OsintDataSource{OsintDataSourceID: 1001, Name: "test_name", Description: "test_desc", MaxScore: 10.0, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Insert",
			input:       &osint.PutOsintDataSourceRequest{ProjectId: 1001, OsintDataSource: &osint.OsintDataSourceForUpsert{Name: "test_name", Description: "test_desc", MaxScore: 10.0}},
			want:        &osint.PutOsintDataSourceResponse{OsintDataSource: &osint.OsintDataSource{OsintDataSourceId: 1001, Name: "test_name", Description: "test_desc", MaxScore: 10.0, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpdResp: &model.OsintDataSource{OsintDataSourceID: 1001, Name: "test_name", Description: "test_desc", MaxScore: 10.0, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid parameter(required project_id)",
			input:   &osint.PutOsintDataSourceRequest{OsintDataSource: &osint.OsintDataSourceForUpsert{Name: "test_name", Description: "test_desc", MaxScore: 10.0, OsintDataSourceId: 1001}},
			wantErr: true,
		},
		{
			name:        "NG DB error(UpsertOsintDataSource)",
			input:       &osint.PutOsintDataSourceRequest{ProjectId: 1001, OsintDataSource: &osint.OsintDataSourceForUpsert{Name: "test_name", Description: "test_desc", MaxScore: 10.0, OsintDataSourceId: 1001}},
			mockUpdResp: nil,
			mockUpdErr:  errors.New("Something error"),
			wantErr:     true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockUpdResp != nil || c.mockUpdErr != nil {
				mockDB.On("UpsertOsintDataSource").Return(c.mockUpdResp, c.mockUpdErr).Once()
			}
			got, err := svc.PutOsintDataSource(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteOsintDataSource(t *testing.T) {
	var ctx context.Context
	mockDB := mockOsintRepository{}
	svc := osintService{repository: &mockDB}
	cases := []struct {
		name     string
		input    *osint.DeleteOsintDataSourceRequest
		wantErr  bool
		mockResp error
	}{
		{
			name:    "OK",
			input:   &osint.DeleteOsintDataSourceRequest{ProjectId: 1001, OsintDataSourceId: 1001},
			wantErr: false,
		},
		{
			name:     "NG DB error",
			input:    &osint.DeleteOsintDataSourceRequest{ProjectId: 1001, OsintDataSourceId: 1001},
			wantErr:  true,
			mockResp: gorm.ErrInvalidSQL,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("DeleteOsintDataSource").Return(c.mockResp).Once()
			_, err := svc.DeleteOsintDataSource(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
		})
	}
}

/**
 * Mock Repository
**/
type mockOsintRepository struct {
	mock.Mock
}

func (m *mockOsintRepository) ListOsint(uint32) (*[]model.Osint, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Osint), args.Error(1)
}
func (m *mockOsintRepository) GetOsint(uint32, uint32) (*model.Osint, error) {
	args := m.Called()
	return args.Get(0).(*model.Osint), args.Error(1)
}
func (m *mockOsintRepository) UpsertOsint(*model.Osint) (*model.Osint, error) {
	args := m.Called()
	return args.Get(0).(*model.Osint), args.Error(1)
}
func (m *mockOsintRepository) DeleteOsint(uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockOsintRepository) ListOsintDataSource(uint32, string) (*[]model.OsintDataSource, error) {
	args := m.Called()
	return args.Get(0).(*[]model.OsintDataSource), args.Error(1)
}
func (m *mockOsintRepository) GetOsintDataSource(uint32, uint32) (*model.OsintDataSource, error) {
	args := m.Called()
	return args.Get(0).(*model.OsintDataSource), args.Error(1)
}
func (m *mockOsintRepository) UpsertOsintDataSource(*model.OsintDataSource) (*model.OsintDataSource, error) {
	args := m.Called()
	return args.Get(0).(*model.OsintDataSource), args.Error(1)
}
func (m *mockOsintRepository) DeleteOsintDataSource(uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockOsintRepository) ListRelOsintDataSource(uint32, uint32, uint32) (*[]model.RelOsintDataSource, error) {
	args := m.Called()
	return args.Get(0).(*[]model.RelOsintDataSource), args.Error(1)
}
func (m *mockOsintRepository) GetRelOsintDataSource(uint32, uint32) (*model.RelOsintDataSource, error) {
	args := m.Called()
	return args.Get(0).(*model.RelOsintDataSource), args.Error(1)
}
func (m *mockOsintRepository) UpsertRelOsintDataSource(*model.RelOsintDataSource) (*model.RelOsintDataSource, error) {
	args := m.Called()
	return args.Get(0).(*model.RelOsintDataSource), args.Error(1)
}
func (m *mockOsintRepository) DeleteRelOsintDataSource(uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockOsintRepository) ListRelOsintDetectWord(uint32, uint32) (*[]model.RelOsintDetectWord, error) {
	args := m.Called()
	return args.Get(0).(*[]model.RelOsintDetectWord), args.Error(1)
}
func (m *mockOsintRepository) GetRelOsintDetectWord(uint32, uint32) (*model.RelOsintDetectWord, error) {
	args := m.Called()
	return args.Get(0).(*model.RelOsintDetectWord), args.Error(1)
}
func (m *mockOsintRepository) UpsertRelOsintDetectWord(*model.RelOsintDetectWord) (*model.RelOsintDetectWord, error) {
	args := m.Called()
	return args.Get(0).(*model.RelOsintDetectWord), args.Error(1)
}
func (m *mockOsintRepository) DeleteRelOsintDetectWord(uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockOsintRepository) ListOsintDetectWord(uint32) (*[]model.OsintDetectWord, error) {
	args := m.Called()
	return args.Get(0).(*[]model.OsintDetectWord), args.Error(1)
}
func (m *mockOsintRepository) GetOsintDetectWord(uint32, uint32) (*model.OsintDetectWord, error) {
	args := m.Called()
	return args.Get(0).(*model.OsintDetectWord), args.Error(1)
}
func (m *mockOsintRepository) UpsertOsintDetectWord(*model.OsintDetectWord) (*model.OsintDetectWord, error) {
	args := m.Called()
	return args.Get(0).(*model.OsintDetectWord), args.Error(1)
}
func (m *mockOsintRepository) DeleteOsintDetectWord(uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockOsintRepository) ListOsintDetectWordFromRelOsintDataSourceID(uint32, uint32) (*[]model.OsintDetectWord, error) {
	args := m.Called()
	return args.Get(0).(*[]model.OsintDetectWord), args.Error(1)
}
func (m *mockOsintRepository) ListAllRelOsintDataSource() (*[]model.RelOsintDataSource, error) {
	args := m.Called()
	return args.Get(0).(*[]model.RelOsintDataSource), args.Error(1)
}
