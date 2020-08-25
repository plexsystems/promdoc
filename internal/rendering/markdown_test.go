package rendering

import (
	"io/ioutil"
	"testing"
)

func TestMarkdown(t *testing.T) {
	expected, err := ioutil.ReadFile("../../examples/markdown/expected.md")
	if err != nil {
		t.Fatal("read expected file")
	}

	actual, err := Render("../../examples/rule.yaml", ".md")
	if err != nil {
		t.Fatal("getting rule groups:", err)
	}

	if string(expected) != actual {
		t.Error("rendered markdown did not match expected output")
	}
}
