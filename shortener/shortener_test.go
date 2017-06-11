package shortener

import "testing"

func TestFnv1a(t *testing.T) {
	if n := fnv1a("node.js"); n != 3096844302 {
		t.Error("Expected 3096844302 got ", n)
	}
}

func TestBase62(t *testing.T) {
	if str := base62(3096844302); str != "dxKcPO" {
		t.Error("Expected dxKcPO got", str)
	}
}

func TestEncode(t *testing.T) {
	id := fnv1a("http://minecraft.project-nemesis.cz/")
	//3860153612
	short := base62(id)

	if enc := Encode("http://minecraft.project-nemesis.cz/"); enc != short {
		t.Error("Expected ", short, " got ", enc)
	}
}
