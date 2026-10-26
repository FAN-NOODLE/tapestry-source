package source

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestResponseConsoleLinkListKOSSyncFlag(t *testing.T) {
	payload, err := json.Marshal(ResponseConsoleLinkList{IsKOSSync: 1})
	if err != nil || !strings.Contains(string(payload), `"is_kos_sync":1`) {
		t.Fatalf("link console list missing is_kos_sync: %s, error=%v", payload, err)
	}
}
