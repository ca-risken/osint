package main

import (
	"fmt"
	"testing"
)

func TestGetScore(t *testing.T) {
	cases := []struct {
		name string
		want float32
	}{
		{
			name: "OK",
			want: 0.1,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getScore()
			if c.want != got {
				t.Fatalf("Unexpected category name: want=%v, got=%v", c.want, got)
			}
		})
	}
}

func TestGetDescription(t *testing.T) {
	cases := []struct {
		name           string
		resourceName   string
		technologyName string
		version        string
		want           string
	}{
		{
			name:           "Exist Version",
			resourceName:   "Hoge",
			technologyName: "Fuga",
			version:        "v0.0.1",
			want:           fmt.Sprintf("%v is using %v. version: %v", "Hoge", "Fuga", "v0.0.1"),
		}, {
			name:           "Not Exist Version",
			resourceName:   "Hoge",
			technologyName: "Fuga",
			version:        "",
			want:           fmt.Sprintf("%v is using %v.", "Hoge", "Fuga"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getDescription(c.resourceName, c.technologyName, c.version)
			if c.want != got {
				t.Fatalf("Unexpected category name: want=%v, got=%v", c.want, got)
			}
		})
	}
}
