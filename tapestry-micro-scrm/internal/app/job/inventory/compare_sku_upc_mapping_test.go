package inventory

import (
	"reflect"
	"testing"
)

func TestDiffUpcCodesReturnsSortedDifferences(t *testing.T) {
	onlyInDB, onlyInPIM, matched := diffUpcCodes([]string{"3", "1", "2"}, []string{"2", "4", "1"})
	if matched {
		t.Fatalf("expected UPC lists to differ")
	}
	if !reflect.DeepEqual(onlyInDB, []string{"3"}) {
		t.Fatalf("unexpected onlyInDB: %#v", onlyInDB)
	}
	if !reflect.DeepEqual(onlyInPIM, []string{"4"}) {
		t.Fatalf("unexpected onlyInPIM: %#v", onlyInPIM)
	}
}

func TestDiffUpcCodesIgnoresOrderAndDuplicateValues(t *testing.T) {
	onlyInDB, onlyInPIM, matched := diffUpcCodes([]string{"2", "1", "1"}, []string{"1", "2"})
	if !matched {
		t.Fatalf("expected UPC lists to match")
	}
	if len(onlyInDB) != 0 || len(onlyInPIM) != 0 {
		t.Fatalf("expected empty differences, got db=%#v pim=%#v", onlyInDB, onlyInPIM)
	}
}

func TestParseManualCompareSkuItems(t *testing.T) {
	got := parseManualCompareSkuItems(" ceu87,CEU87, cew12 ,,abc01 ")
	want := []string{"CEU87", "CEW12", "ABC01"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseManualCompareSkuItems() = %#v, want %#v", got, want)
	}
}
