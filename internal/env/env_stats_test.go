package env

import (
	"testing"
)

func TestStats_EmptyMap(t *testing.T) {
	r := Stats(map[string]string{}, DefaultStatsOptions())
	if r.Total != 0 || r.Empty != 0 || r.Sensitive != 0 {
		t.Fatalf("expected all zeros, got %+v", r)
	}
}

func TestStats_TotalCount(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2", "C": "3"}
	r := Stats(env, DefaultStatsOptions())
	if r.Total != 3 {
		t.Errorf("expected Total=3, got %d", r.Total)
	}
}

func TestStats_EmptyValues(t *testing.T) {
	env := map[string]string{"A": "", "B": "val", "C": ""}
	r := Stats(env, DefaultStatsOptions())
	if r.Empty != 2 {
		t.Errorf("expected Empty=2, got %d", r.Empty)
	}
}

func TestStats_SensitiveKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "secret",
		"API_TOKEN":   "tok",
		"APP_NAME":    "myapp",
	}
	r := Stats(env, DefaultStatsOptions())
	if r.Sensitive != 2 {
		t.Errorf("expected Sensitive=2, got %d", r.Sensitive)
	}
}

func TestStats_UniqueValues(t *testing.T) {
	env := map[string]string{"A": "x", "B": "x", "C": "y"}
	r := Stats(env, DefaultStatsOptions())
	// unique values: "x" and "y" => 2
	if r.Unique != 2 {
		t.Errorf("expected Unique=2, got %d", r.Unique)
	}
}

func TestStats_PrefixCounts(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "prod",
	}
	r := Stats(env, DefaultStatsOptions())
	if r.PrefixCounts["DB"] != 2 {
		t.Errorf("expected DB prefix count=2, got %d", r.PrefixCounts["DB"])
	}
	if r.PrefixCounts["APP"] != 1 {
		t.Errorf("expected APP prefix count=1, got %d", r.PrefixCounts["APP"])
	}
}

func TestStats_TopPrefixes_LimitedByTopN(t *testing.T) {
	env := map[string]string{
		"A_1": "v", "A_2": "v", "A_3": "v",
		"B_1": "v", "B_2": "v",
		"C_1": "v",
		"D_1": "v",
	}
	opts := StatsOptions{TopN: 2}
	r := Stats(env, opts)
	if len(r.TopPrefixes) != 2 {
		t.Errorf("expected 2 top prefixes, got %d", len(r.TopPrefixes))
	}
	if r.TopPrefixes[0].Prefix != "A" {
		t.Errorf("expected top prefix A, got %s", r.TopPrefixes[0].Prefix)
	}
}

func TestStats_TopPrefixes_ZeroTopN_ReturnsAll(t *testing.T) {
	env := map[string]string{"X_1": "a", "Y_1": "b", "Z_1": "c"}
	opts := StatsOptions{TopN: 0}
	r := Stats(env, opts)
	if len(r.TopPrefixes) != 3 {
		t.Errorf("expected 3 prefixes, got %d", len(r.TopPrefixes))
	}
}
