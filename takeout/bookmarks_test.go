package takeout

import (
	"fmt"
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
	dumpFolder(b)
	t.Fail()
}

func dumpFolder(f *BookmarkFolder) {
	fmt.Printf("%s\nAdd Date: %s\nLast Modified: %s\n",
		f.Title, f.AddDate, f.LastModified)
	fmt.Println()
	for i := range f.Entries {
		switch e := f.Entries[i].(type) {
		case *BookmarkFolder:
			dumpFolder(e)
		case *Bookmark:
			dumpBookmark(e)
		}
	}
	fmt.Println("----")
}

func dumpBookmark(b *Bookmark) {
	fmt.Printf("%s %s %s %s\n", b.AddDate, b.Title, b.URL, b.IconURI)
}
