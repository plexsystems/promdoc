package generate

import (
	"io/ioutil"
	"testing"
)

func TestMarkdown(t *testing.T) {
	expectedBytes, err := ioutil.ReadFile("../examples/markdown/expected.md")
	if err != nil {
		t.Fatal("read:", err)
	}

	actual, err := Generate("../examples", ".md")
	if err != nil {
		t.Fatal("generate:", err)
	}

	expected := string(expectedBytes)
	if expected != actual {
		t.Errorf("Unexpected Markdown. expected %v, actual %v", expected, actual)
	}
}
