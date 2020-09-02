package rendering

import (
	"io/ioutil"
	"testing"
)

func testMarkdown(t *testing.T, filename string) {
	expected, err := ioutil.ReadFile("../../examples/markdown/expected.md")
	if err != nil {
		t.Fatal("read expected file")
	}

	actual, err := Render(filename, ".md")
	if err != nil {
		t.Fatal("getting rule groups:", err)
	}

	if string(expected) != actual {
		t.Error("rendered markdown did not match expected output")
	}
}

func TestMarkdownWithDescription(t *testing.T) {
	testCSV(t, "../../examples/rule.yaml")
}

func TestMarkdownWithMessage(t *testing.T) {
	testCSV(t, "../../examples/rule2.yaml")
}
