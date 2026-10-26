package source

import (
	"reflect"
	"testing"

	sourceModel "micro-scrm/internal/app/models/source"
)

func TestInitialGroupThirdIDsUsesAllDepartmentsForRootGroup(t *testing.T) {
	got := initialGroupThirdIDs(0, []int{300, 100, 200}, nil)
	want := []int{100, 200, 300}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected root permissions: got %#v want %#v", got, want)
	}
}

func TestInitialGroupThirdIDsCopiesDirectParentPermissions(t *testing.T) {
	parentRelations := []sourceModel.GroupRelationModel{
		{ThirdId: 301},
		{ThirdId: 101},
		{ThirdId: 301},
	}
	got := initialGroupThirdIDs(99, []int{999}, parentRelations)
	want := []int{101, 301}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected child permissions: got %#v want %#v", got, want)
	}
}

func TestInitialGroupThirdIDsKeepsZeroPermissionParentEmpty(t *testing.T) {
	got := initialGroupThirdIDs(99, []int{100, 200}, nil)
	if len(got) != 0 {
		t.Fatalf("expected no inherited permissions, got %#v", got)
	}
}
