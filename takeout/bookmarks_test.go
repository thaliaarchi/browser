package takeout

import (
	"encoding/json"
	"os"
	"testing"
)

func TestBookmarks(t *testing.T) {
	f, err := os.Open("Bookmarks.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	b, err := ParseNetscapeBookmarks(f)
	if err != nil {
		t.Fatal(err)
	}
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	if err := e.Encode(b); err != nil {
		t.Error(err)
	}
	t.Fail()
}
