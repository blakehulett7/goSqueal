package goSqueal

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestGetTableFields(t *testing.T) {
	defer exec.Command("rm", "database.db").Run()
	CheckForTable("users_test")
	defer DropTable("users_test")
	tests := map[string]struct {
		tableName string
		want      []string
	}{
		"simple": {tableName: "users_test", want: []string{"id", "username", "refresh_token"}},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := GetTableFields(test.tableName)
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("test failed: wanted %v, got: %v", test.want, got)
			}
		})
	}
}

func TestCreateGetTableEntry(t *testing.T) {
	defer exec.Command("rm", "database.db").Run()
	CheckForTable("users_test")
	defer DropTable("users_test")
	tests := map[string]struct {
		tableName string
		params    map[string]string
		want      map[string]string
	}{
		"simple": {
			tableName: "users_test",
			params: map[string]string{
				"id":            "1",
				"username":      "bhulett",
				"refresh_token": "asdf",
			},
			want: map[string]string{"id": "1", "username": "bhulett", "refresh_token": "asdf"},
		},
		"badKey": {
			tableName: "users_test",
			params: map[string]string{
				"id":            "2",
				"name":          "bhulett",
				"refresh_token": "jkl;",
			},
			want: map[string]string{"id": "2", "username": "", "refresh_token": "jkl;"},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			CreateTableEntry(test.tableName, test.params)
			got := GetTableEntry(test.tableName, test.params["id"])
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("test failed: wanted %v, got: %v", test.want, got)
			}
		})
	}
}

func TestDeleteEntry(t *testing.T) {
	defer exec.Command("rm", "database.db").Run()
	CheckForTable("users_test")
	defer DropTable("users_test")
	tests := map[string]struct {
		tableName string
		id        string
		want      map[string]string
	}{
		"simple": {
			"users_test",
			"1",
			nil,
		},
	}
	CreateTableEntry("users_test", map[string]string{"id": "1", "username": "bhulett", "refresh_token": "asdf"})
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			DeleteTableEntry(test.tableName, test.id)
			got := GetTableEntry("users_test", "1")
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("test failed: wanted %v, got: %v", test.want, got)
			}
		})
	}
}
