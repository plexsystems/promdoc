package rendering

import (
	"io/ioutil"
	"testing"
)

func TestMarkdown(t *testing.T) {
	expectedBytes, err := ioutil.ReadFile("../../examples/markdown/expected.md")
	if err != nil {
		t.Fatal("read:", err)
	}

	actual, err := Render("../../examples", ".md")
	if err != nil {
		t.Fatal("render:", err)
	}

	expected := string(expectedBytes)
	if expected != actual {
		t.Errorf("Unexpected markdown. expected %v, actual %v", expected, actual)
	}
}
