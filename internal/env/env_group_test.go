package env

import (
	"testing"
)

func TestGroupMap_NoKeys_ReturnsEmpty(t *testing.T) {
	result := GroupMap(map[string]string{}, DefaultGroupOptions())
	if len(result) != 0 {
		t.Fatalf("expected 0 groups, got %d", len(result))
	}
}

func TestGroupMap_SinglePrefix(t *testing.T) {
	m := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"DB_NAME": "mydb",
	}
	result := GroupMap(m, DefaultGroupOptions())
	if len(result) != 1 {
		t.Fatalf("expected 1 group, got %d", len(result))
	}
	if result[0].Label != "DB" {
		t.Errorf("expected label DB, got %s", result[0].Label)
	}
	if len(result[0].Keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(result[0].Keys))
	}
}

func TestGroupMap_MultiplePrefix_SortedByLabel(t *testing.T) {
	m := map[string]string{
		"REDIS_HOST": "localhost",
		"DB_HOST":    "localhost",
		"APP_PORT":   "8080",
	}
	result := GroupMap(m, DefaultGroupOptions())
	if len(result) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(result))
	}
	if result[0].Label != "APP" || result[1].Label != "DB" || result[2].Label != "REDIS" {
		t.Errorf("unexpected label order: %s %s %s", result[0].Label, result[1].Label, result[2].Label)
	}
}

func TestGroupMap_UngroupedIncluded(t *testing.T) {
	m := map[string]string{
		"HOSTNAME": "box1",
		"DB_HOST":  "localhost",
	}
	opts := DefaultGroupOptions()
	opts.IncludeUngrouped = true
	result := GroupMap(m, opts)
	labels := map[string]bool{}
	for _, g := range result {
		labels[g.Label] = true
	}
	if !labels["OTHER"] {
		t.Error("expected OTHER group for keys without delimiter")
	}
}

func TestGroupMap_UngroupedExcluded(t *testing.T) {
	m := map[string]string{
		"HOSTNAME": "box1",
		"DB_HOST":  "localhost",
	}
	opts := DefaultGroupOptions()
	opts.IncludeUngrouped = false
	result := GroupMap(m, opts)
	for _, g := range result {
		if g.Label == "OTHER" {
			t.Error("did not expect OTHER group when IncludeUngrouped=false")
		}
	}
}

func TestGroupMap_KeysSortedWithinGroup(t *testing.T) {
	m := map[string]string{
		"DB_PORT": "5432",
		"DB_HOST": "localhost",
		"DB_NAME": "mydb",
	}
	result := GroupMap(m, DefaultGroupOptions())
	if len(result) != 1 {
		t.Fatalf("expected 1 group, got %d", len(result))
	}
	keys := result[0].Keys
	if keys[0] != "DB_HOST" || keys[1] != "DB_NAME" || keys[2] != "DB_PORT" {
		t.Errorf("keys not sorted: %v", keys)
	}
}

func TestGroupMap_CustomDelimiter(t *testing.T) {
	m := map[string]string{
		"DB.HOST": "localhost",
		"DB.PORT": "5432",
	}
	opts := DefaultGroupOptions()
	opts.Delimiter = "."
	result := GroupMap(m, opts)
	if len(result) != 1 || result[0].Label != "DB" {
		t.Errorf("expected group DB with dot delimiter, got %+v", result)
	}
}
