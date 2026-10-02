package config

import (
	"reflect"
	"testing"
)

func TestDiffArrays(t *testing.T) {
	oldA := []string{"a", "b", "c"}
	newA := []string{"b", "c", "d"}
	added, removed := DiffArrays(oldA, newA)
	if !reflect.DeepEqual(added, []string{"d"}) {
		t.Fatalf("added want [d], got %v", added)
	}
	if !reflect.DeepEqual(removed, []string{"a"}) {
		t.Fatalf("removed want [a], got %v", removed)
	}

	// 空数组
	added, removed = DiffArrays(nil, []string{"x"})
	if !reflect.DeepEqual(added, []string{"x"}) || len(removed) != 0 {
		t.Fatalf("empty old case: added=%v removed=%v", added, removed)
	}
	added, removed = DiffArrays([]string{"x"}, nil)
	if len(added) != 0 || !reflect.DeepEqual(removed, []string{"x"}) {
		t.Fatalf("empty new case: added=%v removed=%v", added, removed)
	}
}
